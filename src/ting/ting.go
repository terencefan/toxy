package ting

import (
	"errors"
	"fmt"
	"toxy"
	"xlog"

	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/thrift"
	. "github.com/stdrickforce/thriftgo/transport"
)

func Run(conf *toxy.Config) (err error) {
	var service *toxy.ServiceConfig

	defer func() {
		if err != nil {
			xlog.Error(fmt.Sprintf("ping service %s failed:", service.Name))
		}
	}()

	var (
		trans Transport
		proto Protocol
	)

	for _, service = range conf.Services {
		tf, tw, pf, err := parse_config(service)
		if err != nil {
			return err
		}

		trans = tf.GetTransport()
		trans = tw.GetTransport(trans)
		proto = pf.GetProtocol(trans)
		if err := ping(proto); err != nil {
			return err
		}
	}
	return
}

func parse_config(service *toxy.ServiceConfig) (
	tf TransportFactory,
	tw TransportWrapper,
	pf ProtocolFactory,
	err error,
) {
	tf = get_transport_factory(service.Transport, service.Addr)
	tw = get_transport_wrapper(service.Wrapper)
	pf = get_protocol_factory(service.Protocol)
	return
}

func ping(proto Protocol) (err error) {
	if err = proto.GetTransport().Open(); err != nil {
		return
	}
	defer proto.GetTransport().Close()

	if err = proto.WriteMessageBegin("RevenueOrder:ping", T_CALL, 0); err != nil {
		return
	}
	if err = proto.WriteStructBegin("whatever"); err != nil {
		return
	}
	if err = proto.WriteFieldStop(); err != nil {
		return
	}
	if err = proto.WriteStructEnd(); err != nil {
		return
	}
	if err = proto.WriteMessageEnd(); err != nil {
		return
	}
	if err = proto.GetTransport().Flush(); err != nil {
		return
	}
	_, mtype, _, err := proto.ReadMessageBegin()
	if err != nil {
		return
	} else if mtype != T_REPLY {
		// TODO parse replied exception message.
		err = errors.New("reply type mismatch!")
	}
	return
}

func get_transport_factory(tfn, addr string) (tf TransportFactory) {
	switch tfn {
	case "socket":
		tf = NewTSocketFactory(addr)
	case "unix_socket":
		tf = NewTUnixSocketFactory(addr)
	case "http":
		tf = NewTHttpTransportFactory(addr)
	default:
		panic("transport factory invalid!")
	}
	return
}

func get_transport_wrapper(twn string) (tw TransportWrapper) {
	switch twn {
	case "":
		tw = TTransportWrapper
	case "buffered":
		tw = NewTBufferedTransportFactory(4096, 4096)
	case "framed":
		tw = NewTFramedTransportFactory(false, true)
	default:
		panic("transport wrapper invalid!")
	}
	return
}

func get_protocol_factory(pfn string) (pf ProtocolFactory) {
	switch pfn {
	case "binary":
		pf = NewTBinaryProtocolFactory(true, true)
	default:
		panic("protocol factory invalid!")
	}
	return
}
