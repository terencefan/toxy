package main

import (
	"io"
	"net"
	"xhandler"
	"xlog"
	"xprocessor"
	"xprotocol"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	ini "gopkg.in/ini.v1"
)

type Handler struct {
	name string
	h    xhandler.Handler
}

type Gateway struct {
	socket_addr string
	http_addr   string
	processor   xprocessor.Processor
}

func (self *Gateway) process(conn net.Conn) {
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
	}()

	self.processor.Process(conn)
}

func (self *Gateway) Serve() {
	laddr, err := net.ResolveTCPAddr("tcp", self.socket_addr)
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		panic(err)
	}
	xlog.Info("toxy is listening on %s", self.socket_addr)

	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go self.process(conn)
	}
}

func (self *Gateway) ServeHttp() {
}

func (self *Gateway) InitMetric(section *ini.Section) (err error) {
	return
}

func (self *Gateway) InitServices(sections []*ini.Section) (err error) {
	for _, section := range sections {
		name := section.Name()[8:]

		handler, err := xhandler.NewHandler(section)
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

func (self *Gateway) InitProcessor(section *ini.Section) (err error) {
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
	pf := xprotocol.NewTBinaryProtocolFactory(true, true)
	switch ptype {
	case "default":
		fallthrough
	case "single":
		self.processor = xprocessor.NewProcessor(pf)
	case "multiplexed":
		self.processor = xprocessor.NewMultiplexedProcessor(pf)
	default:
		panic("processor type must be one of: [single, multiplexed]")
	}
	return
}

func (self *Gateway) LoadConfig(filepath string) (err error) {
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

func MakeGateway() *Gateway {
	return &Gateway{}
}

var (
	config = kingpin.Flag("config", "Config file.").Short('c').Required().String()
)

func main() {
	// kingpin.Parse()
	var gateway = MakeGateway()
	if err := gateway.LoadConfig("toxy.ini"); err != nil {
		panic(err)
	}
	gateway.Serve()
}
