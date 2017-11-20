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

func (self *SingleProcessor) Process(conn net.Conn) {
	itrans := NewTSocketConn(conn)
	protocol := self.pf.NewProtocol(itrans)
	m := NewMessenger(protocol)

	// Get a server protocol
	if self.handler == nil {
		panic(errors.New("no handler has been set"))
	}
	otrans, err := self.handler.GetTransport()
	if err != nil {
		panic(err)
	}
	oprot := NewTBinaryProtocol(otrans, true, true)
	m.SetOutputProtocol(oprot)

	// close transports after process finished.
	defer itrans.Close()
	defer m.DelOutputProtocol()

	for {
		// TODO support oneway request.
		name, seqid := read_header(m)

		if name == "ping" {
			// reply_shutdown(m, name, seqid)
			fast_reply(m, seqid)
			continue
		}

		if self.shutdown {
			reply_shutdown(m, name, seqid)
			break
		} else {
			reply(m, name, seqid)
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
