package main

import (
	"os"

	"github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/thrift"
	"github.com/stdrickforce/thriftgo/transport"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func ping(proto protocol.Protocol) byte {
	if err := proto.WriteMessageBegin("RevenueOrder:ping", T_CALL, 0); err != nil {
		panic(err)
	}
	if err := proto.WriteStructBegin("whatever"); err != nil {
		panic(err)
	}
	if err := proto.WriteFieldStop(); err != nil {
		panic(err)
	}
	if err := proto.WriteStructEnd(); err != nil {
		panic(err)
	}
	if err := proto.WriteMessageEnd(); err != nil {
		panic(err)
	}
	if err := proto.Flush(); err != nil {
		panic(err)
	}
	_, mtype, _, err := proto.ReadMessageBegin()
	if err != nil {
		panic(err)
	}
	return mtype
}

var (
	addr = kingpin.Flag("addr", "address").Short('a').Default(":6000").String()
	// Transport Factory Name
	tfn = kingpin.Flag("transport_factory", "TransportFactory").Short('t').Default("socket").String()
	// Transport Wrapper Name
	twn = kingpin.Flag("transport_wrapper", "TransportWrapper").Short('w').Default("default").String()
	// Protocol Factory Name
	pfn = kingpin.Flag("protocol_factory", "ProtocolFactory").Short('p').Default("binary").String()

	// Request Path (for http transport)
	path = kingpin.Flag("path", "http request path").Default("/").String()
)

func get_transport_factory() (tf transport.TransportFactory) {
	switch *tfn {
	case "socket":
		tf = transport.NewTSocketFactory(*addr)
	case "unix_socket":
		tf = transport.NewTUnixSocketFactory(*addr)
	case "http":
		tf = transport.NewTHttpTransportFactory(*addr, *path)
	default:
		panic("transport factory invalid!")
	}
	return
}

func get_transport_wrapper() (tw transport.TransportWrapper) {
	switch *twn {
	case "default":
		tw = transport.TTransportWrapper
	case "buffered":
		tw = transport.NewTBufferedTransportFactory(4096, 4096)
	case "framed":
		tw = transport.NewTFramedTransportFactory(false, true)
	default:
		panic("transport wrapper invalid!")
	}
	return
}

func get_protocol_factory() (pf protocol.ProtocolFactory) {
	switch *pfn {
	case "binary":
		pf = protocol.NewTBinaryProtocolFactory(true, true)
	default:
		panic("protocol factory invalid!")
	}
	return
}

func main() {
	kingpin.Parse()
	tf := get_transport_factory()
	tw := get_transport_wrapper()
	pf := get_protocol_factory()

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
