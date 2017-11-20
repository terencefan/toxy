package xprocessor

import (
	"net"
	"strings"
	"xhandler"
	"xlog"
	. "xprotocol"
	. "xtransport"
)

type MultiplexedProcessor struct {
	pf       ProtocolFactory
	hmap     map[string]*xhandler.Handler
	shutdown bool
}

func (self *MultiplexedProcessor) Add(
	name string,
	handler *xhandler.Handler,
) (err error) {
	self.hmap[name] = handler
	return
}

func (self *MultiplexedProcessor) parse_name(name string) (
	service, fname string,
) {
	segments := strings.SplitN(name, ":", 2)
	if len(segments) != 2 {
		err := NewProcessorError("fname format mismatch: %s", name)
		panic(err)
	}
	service, fname = segments[0], segments[1]
	return
}

func (self *MultiplexedProcessor) get_protocol(service string) (
	protocol Protocol,
) {
	handler := self.hmap[service]
	if handler == nil {
		err := NewProcessorError("no handler has been set for: %s", service)
		panic(err)
	}

	otrans, err := handler.GetTransport()
	if err != nil {
		panic(err)
	}
	protocol = NewTBinaryProtocol(otrans, true, true)
	return
}

func (self *MultiplexedProcessor) Process(conn net.Conn) {
	itrans := NewTSocketConn(conn)
	protocol := self.pf.NewProtocol(itrans)
	m := NewMessenger(protocol)

	defer itrans.Close()

	for {
		// TODO support oneway request.
		name, seqid := read_header(m)

		service, fname := self.parse_name(name)
		xlog.Debug("%s: %s", service, fname)

		// NOTE fast reply ping requests.
		if name == "ping" {
			fast_reply(m, seqid)
			continue
		}

		oprot := self.get_protocol(service)
		m.SetOutputProtocol(oprot)

		reply(m, fname, seqid)

		m.DelOutputProtocol()

		// if self.shutdown {
		// 	m.Reverse()
		// 	err := errors.New("server is going away ~!")
		// 	write_header(m, fname, T_EXCEPTION, seqid)
		// 	write_body_error(m, err)
		// 	return
		// }
	}
}

func (self *MultiplexedProcessor) Shutdown() (err error) {
	return
}

func NewMultiplexedProcessor(pf ProtocolFactory) *MultiplexedProcessor {
	return &MultiplexedProcessor{
		pf:       pf,
		hmap:     make(map[string]*xhandler.Handler),
		shutdown: false,
	}
}
