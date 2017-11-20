package xprocessor

import (
	"net"
	"strings"
	"xhandler"
	"xlog"
	. "xprotocol"
	. "xthrift"
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
		} else {
			oprot := self.get_protocol(service)
			m.SetOutputProtocol(oprot)
			reply(m, fname, seqid)
			m.DelOutputProtocol()
		}

		// if self.shutdown {
		// 	m.Reverse()
		// 	err := errors.New("server is going away ~!")
		// 	write_header(m, fname, T_EXCEPTION, seqid)
		// 	write_body_error(m, err)
		// 	return
		// }
	}

}

func NewMultiplexedProcessor(pf ProtocolFactory) *MultiplexedProcessor {
	return &MultiplexedProcessor{
		pf:       pf,
		hmap:     make(map[string]*xhandler.Handler),
		shutdown: false,
	}
}

func read_header(m *Messenger) (name string, seqid int32) {
	xlog.Debug("read message header")
	name, mtype, seqid, err := m.ReadMessageBegin()
	if err != nil {
		panic(err)
	} else if mtype == T_ONEWAY {
		// TODO reply exception "doesn't support oneway request yet."
	} else if mtype != T_CALL {
		// TODO raise exception.
	} else {
		return
	}
	return
}

func write_header(m *Messenger, name string, mtype byte, seqid int32) {
	xlog.Debug("write message header")
	if err := m.WriteMessageBegin(name, mtype, seqid); err != nil {
		panic(err)
	}
}

func forward_header(m *Messenger) {
	xlog.Debug("forward message header")
	if err := m.ForwardMessageBegin(); err != nil {
		panic(err)
	}
}

func forward_body(m *Messenger) {
	xlog.Debug("forward message body")
	if err := m.Forward(T_STRUCT); err != nil {
		panic(err)
	}
	if err := m.ForwardMessageEnd(); err != nil {
		panic(err)
	}
}

func fast_reply(m *Messenger, seqid int32) {
	if err := m.FastReply(seqid); err != nil {
		panic(err)
	}
}

func reply(m *Messenger, name string, seqid int32) {
	write_header(m, name, T_CALL, seqid)
	forward_body(m)
	m.Reverse()
	forward_header(m)
	forward_body(m)
	m.Reverse()
}
