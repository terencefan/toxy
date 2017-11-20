package xprocessor

import (
	"fmt"
	"net"
	"xhandler"
)

type Processor interface {
	Add(string, *xhandler.Handler) error
	Process(net.Conn)
	Shutdown() error
}

type ProcessorError struct {
	Message string
}

func (e ProcessorError) Error() string {
	return e.Message
}

func NewProcessorError(format string, args ...interface{}) ProcessorError {
	return ProcessorError{
		Message: fmt.Sprintf(format, args...),
	}
}
