package toxy

import (
	"fmt"
	"strings"

	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/thrift"
)

type MultiplexedProcessor struct {
	hmap map[string]*Handler
}

func (self *MultiplexedProcessor) Add(
	name string,
	handler *Handler,
) (err error) {
	self.hmap[name] = handler
	return
}

func (self *MultiplexedProcessor) Parse(name string) (
	service, fname string, err error,
) {
	segments := strings.SplitN(name, ":", 2)
	if len(segments) != 2 {
		err = NewTApplicationException(
			fmt.Sprintf(
				"Service name not found in message name: %s. "+
					"Did you forget to use a TMultiplexProtocol in your client?",
				name,
			),
			ExceptionWrongMethodName,
		)
	} else {
		service, fname = segments[0], segments[1]
	}
	return
}

func (self *MultiplexedProcessor) GetProtocol(name string) (Protocol, error) {
	if handler, ok := self.hmap[name]; ok {
		return handler.GetProtocol()
	} else {
		return nil, NewTApplicationException(
			fmt.Sprintf("Service `%s` has not been registered", name),
			ExceptionUnknownMethod,
		)
	}
}

func NewMultiplexedProcessor() *MultiplexedProcessor {
	return &MultiplexedProcessor{
		hmap: make(map[string]*Handler),
	}
}
