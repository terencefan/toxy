package toxy

import "github.com/stdrickforce/thriftgo/protocol"

type Processor interface {
	Add(string, *Handler) error
	Parse(name string) (fname, service string, err error)
	GetProtocol(service string) (protocol.Protocol, error)
}
