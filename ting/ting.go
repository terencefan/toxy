package main

import (
	"os"

	"github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/thrift"
	"github.com/stdrickforce/thriftgo/transport"
)

func ping(proto protocol.Protocol) byte {
	proto.WriteMessageBegin("RevenueOrder:ping", T_CALL, 0)
	proto.WriteStructBegin("whatever")
	proto.WriteFieldStop()
	proto.WriteStructEnd()
	proto.WriteMessageEnd()
	proto.Flush()
	_, mtype, _, err := proto.ReadMessageBegin()
	if err != nil {
		panic(err)
	}
	return mtype
}

func main() {
	tw := transport.NewTFramedTransportFactory(false, true)
	tf := transport.NewTSocketFactory("localhost:10011")
	pf := protocol.NewTBinaryProtocolFactory(true, true)

	var (
		trans transport.Transport
		proto protocol.Protocol
		err   error
	)

	trans, err = tf.GetTransport()
	if err != nil {
		panic(err)
	}

	trans, err = tw.Wraps(trans)
	if err != nil {
		panic(err)
	}
	proto = pf.NewProtocol(trans)

	mtype := ping(proto)
	if mtype == T_REPLY {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
