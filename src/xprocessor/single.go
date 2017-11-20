package xprocessor

import (
	"errors"
	"net"
	"xhandler"
	. "xprotocol"
	. "xtransport"
)

type SingleProcessor struct {
	name     string
	handler  *xhandler.Handler
	pf       ProtocolFactory
	shutdown bool
}

func (self *SingleProcessor) Add(
	name string,
	handler *xhandler.Handler,
) (err error) {
	self.name = name
	self.handler = handler
	return
}

func (self *SingleProcessor) handle(m *Messenger) bool {
	name, seqid := read_header(m)

	if name == "ping" {
		// reply_shutdown(m, name, seqid)
		fast_reply(m, seqid)
		return true
	}

	if self.shutdown {
		reply_shutdown(m, name, seqid)
		return false
	} else {
		reply(m, name, seqid)
		return true
	}

}

func (self *SingleProcessor) get_protocol() Protocol {
	// Get a server protocol
	if self.handler == nil {
		panic(errors.New("no handler has been set"))
	}
	otrans, err := self.handler.GetTransport()
	if err != nil {
		panic(err)
	}
	return NewTBinaryProtocol(otrans, true, true)
}

func (self *SingleProcessor) Process(conn net.Conn) {
	itrans := NewTSocketConn(conn)
	defer itrans.Close()

	protocol := self.pf.NewProtocol(itrans)
	m := NewMessenger(protocol)

	oprot := self.get_protocol()
	m.SetOutputProtocol(oprot)
	defer m.DelOutputProtocol()

	for {
		if ok := self.handle(m); !ok {
			break
		}
	}
}

func (self *SingleProcessor) Shutdown() (err error) {
	self.shutdown = true
	return
}

func NewProcessor(pf ProtocolFactory) *SingleProcessor {
	return &SingleProcessor{
		pf:       pf,
		shutdown: false,
	}
}
