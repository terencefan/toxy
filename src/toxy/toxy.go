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

	raven "github.com/getsentry/raven-go"
	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/thrift"
	. "github.com/stdrickforce/thriftgo/transport"
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
	if err == io.EOF {
		// NOTE reset by client or server ?
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
				xlog.Error(err.Error())
				return
			}
			if err = send_exception(proto, ae); err != nil {
				xlog.Error(err.Error())
				return
			}
			loop = ae.Type != ExceptionShutdown
		}
		xlog.Warning(ae.Error())
		raven.CaptureError(ae, nil)
	} else {
		xmetric.Count("toxy", "error.unexpected", 1)
		raven.CaptureError(err, nil)
		xlog.Error(err.Error())
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
	xlog.Info("receive request")

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

	// fast reply ping requests.
	if name == "ping" && self.fast_reply {
		defer xlog.Info("fast reply ping request")
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
	xlog.Info("send response")
	return
}

func (self *Toxy) loop(conn net.Conn) {
	self.wg.Add(1)

	defer conn.Close()
	defer self.wg.Done()

	defer xmetric.Count("toxy", "connection.closed", 1)
	defer xlog.Debug("connection closed")

	xlog.Info("receive connection")
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

	xlog.Info("toxy is listening on " + self.socket_addr)

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
			} else {
				xlog.Warning(err.Error())
				continue
			}
		}
		go self.loop(conn)
	}
	self.wg.Wait()
}

func (self *Toxy) InitMetric(conf *MetricConfig) (err error) {
	xmetric.InitBufferedStatsd(
		xmetric.Address(conf.Addr),
		xmetric.Prefix(conf.Prefix),
		xmetric.FlushPeriod(time.Second),
		xmetric.MaxBufferSize(1450),
		xmetric.MaxQueueSize(128),
	)
	return
}

func (self *Toxy) InitSentry(conf *SentryConfig) (err error) {
	var context = make(map[string]string)

	// hostname
	if hostname, err := os.Hostname(); err != nil {
		return err
	} else {
		context["hostname"] = hostname
	}

	// appname
	context["app"] = "bank"
	return raven.SetDSN(conf.Dsn)
}

func (self *Toxy) InitProcessor(conf *ProxyConfig) (err error) {
	self.socket_addr = conf.Addr
	switch conf.Processor {
	case "default":
		fallthrough
	case "single":
		self.processor = NewProcessor()
	case "multiplexed":
		self.processor = NewMultiplexedProcessor()
	default:
		panic("processor type must be one of: [single, multiplexed]")
	}
	return
}

func (self *Toxy) AddService(conf *ServiceConfig) (err error) {
	handler, err := NewHandler(conf.Name, conf)
	if err != nil {
		return err
	}
	err = self.processor.Add(conf.Name, handler)
	if err != nil {
		return err
	}
	return
}

func (self *Toxy) FastReply() {
	self.fast_reply = true
}

func (self *Toxy) Init(conf *Config) (err error) {
	// initialize processor
	if err = self.InitProcessor(conf.Proxy); err != nil {
		return
	}
	xlog.Info(fmt.Sprintf("proxy is running in %s mode", conf.Proxy.Processor))

	// initialize metric
	if conf.Metric != nil {
		if err = self.InitMetric(conf.Metric); err != nil {
			return
		}
		xlog.Info("metric module has been initialized")
	}

	// initialize sentry
	if conf.Sentry != nil {
		if err = self.InitSentry(conf.Sentry); err != nil {
			return
		}
		xlog.Info("sentry module has been initialized")
	}

	// initialize services
	for _, sc := range conf.Services {
		if err = self.AddService(sc); err != nil {
			return
		}
		xlog.Info(fmt.Sprintf("service %s has been registered", sc.Name))
	}
	return
}

func NewToxy(conf *Config) (toxy *Toxy) {
	toxy = &Toxy{
		wg: new(sync.WaitGroup),
	}
	if err := toxy.Init(conf); err != nil {
		panic(err)
	}
	return
}
