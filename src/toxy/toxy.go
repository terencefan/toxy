package toxy

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"xlog"
	"xmetric"

	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/thrift"
	. "github.com/stdrickforce/thriftgo/transport"
	ini "gopkg.in/ini.v1"
)

var shutdown int32 = 0

const (
	ExceptionServiceUnavailable = 100
	ExceptionShutdown           = 101
)

type Toxy struct {
	socket_addr string
	http_addr   string
	processor   Processor
	wg          *sync.WaitGroup
	fast_reply  bool
}

func send_exception(proto Protocol, ae *TApplicationException) (err error) {
	if err = proto.WriteMessageBegin("unknown", T_EXCEPTION, 0); err != nil {
		return
	}
	if err = WriteTApplicationException(proto, ae); err != nil {
		return
	}
	if err = proto.WriteMessageEnd(); err != nil {
		return
	}
	if err = proto.GetTransport().Flush(); err != nil {
		return
	}
	return
}

func skip_message_body(proto Protocol) (err error) {
	if err = proto.Skip(T_STRUCT); err != nil {
		return
	}
	if err = proto.ReadMessageEnd(); err != nil {
		return
	}
	return
}

func handle_err(proto Protocol, err error) (loop bool) {
	fmt.Println(err)
	if err == io.EOF {
		xmetric.Count("toxy", "connection.closed", 1)
		// NOTE reset by client or server ?
		xlog.Debug("connection reset by peer")
		return
	} else if ae, ok := err.(*TApplicationException); ok {
		switch ae.Type {
		case ExceptionUnknownMethod:
			fallthrough
		case ExceptionInvalidMessageType:
			fallthrough
		case ExceptionWrongMethodName:
			fallthrough
		case ExceptionServiceUnavailable:
			fallthrough
		case ExceptionShutdown:
			if err = skip_message_body(proto); err != nil {
				return false
			}
			if err = send_exception(proto, ae); err != nil {
				return false
			}
			loop = ae.Type != ExceptionShutdown
		}
	} else {
		xmetric.Count("toxy", "error.unexpected", 1)
		xlog.Warning("unexpected err found: %s", err)
	}
	return
}

func (self *Toxy) process(iprot Protocol) (err error) {
	var (
		s_time int64
		key    string
	)

	s_time = time.Now().UnixNano()

	// read message begin from input protocol
	name, mtype, seqid, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	} else if mtype != T_CALL {
		return NewTApplicationException(
			fmt.Sprintf("invalid message type: %d", mtype),
			ExceptionInvalidMessageType,
		)
	}

	// graceful shutdown.
	if atomic.LoadInt32(&shutdown) > 0 {
		return NewTApplicationException(
			"Toxy: proxy is shutting down.",
			ExceptionShutdown,
		)
	}

	// metrics
	// TODO err metric ?
	key = strings.Replace(name, ":", ".", -1)
	xmetric.Count("toxy", key, 1)
	defer func() {
		delta := int((time.Now().UnixNano() - s_time) / 1000000)
		xmetric.Timing("toxy", key, delta)
	}()

	// get output protocol and function name.
	service, name, err := self.processor.Parse(name)
	if err != nil {
		return
	}
	fmt.Println(service, name)

	// fast reply ping requests.
	if name == "ping" && !self.fast_reply {
		return messenger.FastReply(iprot, "ping", seqid)
	}

	// prepare protocols.
	oprot, err := self.processor.GetProtocol(service)
	if err != nil {
		return NewTApplicationException(
			"Toxy: backend server temporarily unavailable: "+err.Error(),
			ExceptionServiceUnavailable,
		)
	}
	defer oprot.GetTransport().Close()

	// forword messages.
	siprot := NewStoredProtocol(iprot, name, T_CALL, seqid)
	if err = messenger.ForwardMessage(siprot, oprot); err != nil {
		return
	}
	if err = messenger.ForwardMessage(oprot, iprot); err != nil {
		return
	}
	return
}

func (self *Toxy) loop(conn net.Conn) {
	self.wg.Add(1)
	defer self.wg.Done()

	remote := conn.RemoteAddr().String()
	xlog.Info("[%s] receive connection", remote)
	xmetric.Count("toxy", "connection.established", 1)

	var (
		trans Transport
		proto Protocol
	)
	trans = NewTSocketConn(conn)
	trans = NewTBufferedTransport(trans)
	proto = NewTBinaryProtocol(trans, true, true)
	defer trans.Close()

	for {
		if err := self.process(proto); err != nil {
			if ok := handle_err(proto, err); ok {
				continue
			} else {
				break
			}
		}
	}
}

func (self *Toxy) Serve() {
	xmetric.Count("toxy", "restart", 1)

	laddr, err := net.ResolveTCPAddr("tcp", self.socket_addr)
	if err != nil {
		panic(err)
	}

	ln, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		panic(err)
	}

	xlog.Info("toxy is listening on %s", self.socket_addr)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		xlog.Warning("Received INT/TERM signal, stopping...")
		atomic.StoreInt32(&shutdown, 1)
		if ln != nil {
			ln.Close()
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if atomic.LoadInt32(&shutdown) > 0 {
				break
			}
			xlog.Warning(err.Error())
			continue
		}
		go self.loop(conn)
	}
	self.wg.Wait()
}

func (self *Toxy) InitMetric(section *ini.Section) (err error) {
	var addr = "0.0.0.0:8125"
	if section.HasKey("addr") {
		addr = section.Key("addr").String()
	}

	var prefix = ""
	if section.HasKey("prefix") {
		prefix = section.Key("prefix").String()
	}

	xmetric.InitBufferedStatsd(
		xmetric.Address(addr),
		xmetric.Prefix(prefix),
		xmetric.FlushPeriod(time.Second),
		xmetric.MaxBufferSize(1450),
		xmetric.MaxQueueSize(128),
	)
	return
}

func (self *Toxy) InitServices(sections []*ini.Section) (err error) {
	for _, section := range sections {
		name := section.Name()[8:]

		handler, err := NewHandler(name, section)
		if err != nil {
			return err
		}

		err = self.processor.Add(name, handler)
		if err != nil {
			return err
		}
	}
	return
}

func (self *Toxy) InitProcessor(section *ini.Section) (err error) {
	if section.HasKey("addr") {
		self.socket_addr = section.Key("addr").String()
	} else {
		self.socket_addr = "0.0.0.0:6000"
	}

	ptype := "default"
	if section.HasKey("processor") {
		ptype = section.Key("processor").String()
	}

	switch ptype {
	case "default":
		fallthrough
	// case "single":
	// 	self.processor = NewProcessor()
	case "multiplexed":
		self.processor = NewMultiplexedProcessor()
	default:
		panic("processor type must be one of: [single, multiplexed]")
	}
	return
}

func (self *Toxy) FastReply() {
	self.fast_reply = true
}

func (self *Toxy) LoadConfig(filepath string) (err error) {
	var f *ini.File
	var section *ini.Section

	// load config file
	if f, err = ini.Load(filepath); err != nil {
		return err
	}

	// initialize metric client
	if section, err = f.GetSection("metric"); err != nil {
		return
	}
	if err = self.InitMetric(section); err != nil {
		return
	}

	// initialize socketserver & processor
	if section, err = f.GetSection("socketserver"); err != nil {
		return
	}
	if err = self.InitProcessor(section); err != nil {
		return
	}

	// TODO init httpserver.

	// TODO init downgrade.

	// initialize backend services.
	if err = self.InitServices(f.ChildSections("service")); err != nil {
		return
	}
	return
}

func NewToxy() *Toxy {
	return &Toxy{
		wg: new(sync.WaitGroup),
	}
}
