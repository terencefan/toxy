package xception

import (
	"fmt"
	. "xprotocol"
	. "xthrift"
)

const (
	ExceptionUnknown            int32 = 0
	ExceptionUnknownMethod            = 1
	ExceptionInvalidMessageType       = 2
	ExceptionWrongMethodName          = 3
	ExceptionBadSequenceID            = 4
	ExceptionMissingResult            = 5
	ExceptionInternalError            = 6
	ExceptionProtocolError            = 7
	// ExceptionInvalidTransform            = 8
	// ExceptionInvalidProtocol             = 9
	// ExceptionUnsupportedClientType       = 10

	// custom error code
	ExceptionShutdown = 100
)

type TApplicationException struct {
	Message string
	Type    int32
}

func (e *TApplicationException) String() string {
	typeStr := "Unknown Exception"
	switch e.Type {
	case ExceptionUnknownMethod:
		typeStr = "Unknown Method"
	case ExceptionInvalidMessageType:
		typeStr = "Invalid Message Type"
	case ExceptionWrongMethodName:
		typeStr = "Wrong Method Name"
	case ExceptionBadSequenceID:
		typeStr = "Bad Sequence ID"
	case ExceptionMissingResult:
		typeStr = "Missing Result"
	case ExceptionInternalError:
		typeStr = "Internal Error"
	case ExceptionProtocolError:
		typeStr = "Protocol Error"
	}
	return fmt.Sprintf("%s: %s", typeStr, e.Message)
}

func (e *TApplicationException) Write(proto Protocol) (err error) {
	if err = proto.WriteStructBegin("TApplicationException"); err != nil {
		return
	}
	if err = proto.WriteFieldBegin("message", T_STRING, 1); err != nil {
		return
	}
	if err = proto.WriteString(e.Message); err != nil {
		return
	}
	if err = proto.WriteFieldEnd(); err != nil {
		return
	}
	if err = proto.WriteFieldBegin("type", T_I32, 2); err != nil {
		return
	}
	if err = proto.WriteI32(e.Type); err != nil {
		return
	}
	if err = proto.WriteFieldEnd(); err != nil {
		return
	}
	if err = proto.WriteFieldStop(); err != nil {
		return
	}
	if err = proto.WriteStructEnd(); err != nil {
		return
	}
	return
}

func NewTApplicationException(message string, t int32) *TApplicationException {
	return &TApplicationException{
		Message: message,
		Type:    t,
	}
}
