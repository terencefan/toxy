package xtransport

import "io"

type Transport interface {
	io.ReadWriteCloser
	Flush() error
}

type TransportFactory interface {
	GetTransport() (Transport, error)
}

type TransportWrapper interface {
	Wraps(Transport) (Transport, error)
}

type transportWrapper struct {
}

func (self *transportWrapper) Wraps(trans Transport) (Transport, error) {
	return trans, nil
}

var TTransportWrapper = &transportWrapper{}
