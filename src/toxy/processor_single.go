package toxy

import (
	"errors"
	"fmt"
	"net"
	"time"
	"xmetric"

	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/thrift"
	. "github.com/stdrickforce/thriftgo/transport"
)

type SingleProcessor struct {
	name    string
	handler *Handler
	pf      ProtocolFactory
}

func (self *SingleProcessor) Add(
	name string,
	handler *Handler,
) (err error) {
	self.name = name
	self.handler = handler
	return
}

func (self *SingleProcessor) handle(iprot, oprot Protocol) bool {
	if shutdown > 0 {
		fast_reply_shutdown(iprot)
		return false
	}

	s_time := time.Now().UnixNano()

	name, seqid := read_message_begin(iprot)

	key := fmt.Sprintf("%s.%s", self.name, name)
	defer func() {
		delta := int((time.Now().UnixNano() - s_time) / 1000000)
		xmetric.Timing("toxy", key, delta)
	}()
	xmetric.Count("toxy", key, 1)

	// NOTE fast reply ping requests.
	// is it neccessary?
	if name == "ping" {
		fast_reply(iprot, "ping", seqid)
		return true
	}

	reply(NewStoredProtocol(iprot, name, T_CALL, seqid), oprot)
	return true
}

func (self *SingleProcessor) get_protocol() Protocol {
	// Get a server protocol
	if self.handler == nil {
		panic(errors.New("no handler has been set"))
	}
	proto, err := self.handler.GetProtocol()
	if err != nil {
		panic(err)
	}
	return proto
}

func (self *SingleProcessor) Process(conn net.Conn) {
	var itrans Transport
	itrans = NewTSocketConn(conn)
	itrans = NewTBufferedTransport(itrans)
	defer itrans.Close()

	iprot := self.pf.GetProtocol(itrans)
	oprot := self.get_protocol()

	for {
		if ok := self.handle(iprot, oprot); !ok {
			break
		}
	}
}

func NewProcessor(pf ProtocolFactory) *SingleProcessor {
	return &SingleProcessor{
		pf: pf,
	}
}
