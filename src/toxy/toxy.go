package toxy

import (
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"xlog"

	. "github.com/stdrickforce/thriftgo/protocol"
	ini "gopkg.in/ini.v1"
)

type Toxy struct {
	socket_addr string
	http_addr   string
	processor   Processor
	shutdown    bool
	wg          *sync.WaitGroup
}

func (self *Toxy) process(conn net.Conn) {
	remote := conn.RemoteAddr().String()
	xlog.Info("[%s] receive connection", remote)

	defer func() {
		if err := recover(); err != nil {
			if err != io.EOF {
				xlog.Info("[%s] unexpected err found: %s", remote, err)
			} else {
				xlog.Debug("[%s] connection reset by peer", remote)
			}
		}
		xlog.Info("[%s] connection closed", remote)
		self.wg.Done()
	}()

	self.processor.Process(conn)
}

func (self *Toxy) Serve() {
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
		self.shutdown = true
		self.processor.Shutdown()
		if ln != nil {
			ln.Close()
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if self.shutdown {
				break
			}
			xlog.Warning(err.Error())
			continue
		}
		self.wg.Add(1)
		go self.process(conn)
	}
	self.wg.Wait()
}

func (self *Toxy) InitMetric(section *ini.Section) (err error) {
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

	// TODO multiple protocol support
	pf := NewTBinaryProtocolFactory(true, true)
	switch ptype {
	case "default":
		fallthrough
	case "single":
		self.processor = NewProcessor(pf)
	case "multiplexed":
		self.processor = NewMultiplexedProcessor(pf)
	default:
		panic("processor type must be one of: [single, multiplexed]")
	}
	return
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

	// TODO httpserver.

	// TODO downgrade.

	// initialize backend services.
	if err = self.InitServices(f.ChildSections("service")); err != nil {
		return
	}
	return
}

func NewToxy() *Toxy {
	return &Toxy{
		shutdown: false,
		wg:       new(sync.WaitGroup),
	}
}
