package toxy

import (
	"fmt"

	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/thrift"
)

type SingleProcessor struct {
	name    string
	handler *Handler
}

func (self *SingleProcessor) Add(
	name string,
	handler *Handler,
) (err error) {
	self.name = name
	self.handler = handler
	return
}

func (self *SingleProcessor) GetProtocol(service string) (Protocol, error) {
	// Get a server protocol
	if self.handler == nil {
		return nil, NewTApplicationException(
			fmt.Sprintf("Service has not been specified"),
			ExceptionUnknownMethod,
		)
	}
	return self.handler.GetProtocol()
}

func (self *SingleProcessor) Parse(name string) (
	service, fname string, err error,
) {
	return self.name, name, nil
}

func NewProcessor() *SingleProcessor {
	return &SingleProcessor{}
}
