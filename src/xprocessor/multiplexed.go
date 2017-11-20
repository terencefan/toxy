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

func (self *MultiplexedProcessor) get_protocol(service string) Protocol {
	handler := self.hmap[service]
	if handler == nil {
		err := NewProcessorError("no handler has been set for: %s", service)
		panic(err)
	}

	otrans, err := handler.GetTransport()
	if err != nil {
		panic(err)
	}
	return NewTBinaryProtocol(otrans, true, true)
}

func (self *MultiplexedProcessor) handle(m *Messenger) bool {
	name, seqid := read_header(m)

	service, fname := self.parse_name(name)
	xlog.Debug("%s: %s", service, fname)

	// NOTE fast reply ping requests.
	if fname == "ping" {
		fast_reply(m, seqid)
		return true
	}

	oprot := self.get_protocol(service)
	m.SetOutputProtocol(oprot)
	defer m.DelOutputProtocol()

	if self.shutdown {
		reply_shutdown(m, fname, seqid)
		return false
	} else {
		reply(m, fname, seqid)
		return true
	}
}

func (self *MultiplexedProcessor) Process(conn net.Conn) {
	itrans := NewTSocketConn(conn)
	defer itrans.Close()

	protocol := self.pf.NewProtocol(itrans)
	m := NewMessenger(protocol)

	for {
		if ok := self.handle(m); !ok {
			break
		}
	}
}

func (self *MultiplexedProcessor) Shutdown() (err error) {
	self.shutdown = true
	return
}

func NewMultiplexedProcessor(pf ProtocolFactory) *MultiplexedProcessor {
	return &MultiplexedProcessor{
		pf:       pf,
		hmap:     make(map[string]*xhandler.Handler),
		shutdown: false,
	}
}
