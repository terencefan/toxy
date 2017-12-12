package toxy

import (
	"errors"
	"time"

	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/transport"

	xparser "github.com/stdrickforce/go-thrift/parser"
)

type Handler struct {
	name        string
	tf          TransportFactory
	tw          TransportWrapper
	multiplexed bool
	timeout     time.Duration
}

var p = &xparser.Parser{}

func new_http_tf(conf *ServiceConfig) (tf TransportFactory, err error) {
	tf = NewTHttpTransportFactory(conf.Addr)
	return
}

func new_socket_tf(conf *ServiceConfig) (tf TransportFactory, err error) {
	tf = NewTSocketFactory(conf.Addr)
	return
}

func new_unix_socket_tf(conf *ServiceConfig) (tf TransportFactory, err error) {
	tf = NewTUnixSocketFactory(conf.Addr)
	return
}

func new_buffered_tw(conf *ServiceConfig) (tw TransportWrapper, err error) {
	rbufsize := 4096
	wbufsize := 4096
	tw = NewTBufferedTransportFactory(rbufsize, wbufsize)
	return
}

func new_framed_tw(conf *ServiceConfig) (tw TransportWrapper, err error) {
	rframed, wframed := false, true
	tw = NewTFramedTransportFactory(rframed, wframed)
	return
}

func NewHandler(name string, conf *ServiceConfig) (h *Handler, err error) {
	h = &Handler{
		name:        name,
		tw:          TTransportWrapper,
		multiplexed: false,
		timeout:     time.Second * 30,
	}

	// transport
	switch conf.Transport {
	case "http":
		h.tf, err = new_http_tf(conf)
	case "socket":
		h.tf, err = new_socket_tf(conf)
	case "unix_socket":
		h.tf, err = new_unix_socket_tf(conf)
	case "tls_socket":
		err = errors.New("Not implement error")
	default:
		// TODO more error messages.
		err = errors.New("invalid transport: " + conf.Transport)
	}

	if err != nil {
		return nil, err
	}

	// wrapper
	switch conf.Wrapper {
	case "":
		h.tw = TTransportWrapper
	case "buffered":
		h.tw, err = new_buffered_tw(conf)
	case "framed":
		h.tw, err = new_framed_tw(conf)
	default:
		err = errors.New("invalid transport wrapper: " + conf.Wrapper)
	}

	if err != nil {
		return nil, err
	}

	h.multiplexed = conf.Multiplexed
	h.timeout = time.Millisecond * time.Duration(conf.Timeout)
	return
}

func (h *Handler) GetTransport() (trans Transport, err error) {
	trans = h.tf.GetTransport()
	trans = h.tw.GetTransport(trans)
	return
}

func (h *Handler) GetProtocol() (proto Protocol, err error) {
	var trans Transport
	if trans, err = h.GetTransport(); err != nil {
		return
	}
	trans.SetTimeout(h.timeout)
	if err = trans.Open(); err != nil {
		return
	}
	proto = NewTBinaryProtocol(trans, true, true)
	if h.multiplexed {
		proto = NewMultiplexedProtocol(proto, h.name)
	}
	return
}
