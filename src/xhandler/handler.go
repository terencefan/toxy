package xhandler

import (
	"errors"
	"xthrift"
	. "xtransport"

	"github.com/stdrickforce/go-thrift/parser"
	ini "gopkg.in/ini.v1"
)

type Handler struct {
	tf TransportFactory
	tw TransportWrapper
}

var p = &parser.Parser{}

func NewHandler(section *ini.Section) (h *Handler, err error) {
	h = &Handler{
		tw: TTransportWrapper,
	}

	// transport
	if key, err := section.GetKey("transport"); err != nil {
		return nil, err
	} else {
		switch key.String() {
		case "http":
			// TODO check if addr / path exist in config
			addr := section.Key("addr").String()
			path := section.Key("path").String()
			h.tf = NewTHttpTransportFactory(addr, path)
		default:
			// TODO more error messages.
			err = errors.New("service config error")
			return nil, err
		}
	}

	// thrift
	if key, err := section.GetKey("thrift"); err != nil {
		// skip
	} else {
		filename := key.String()
		_, _, err := xthrift.Parse(filename)
		if err != nil {
			return nil, err
		}
	}
	return
}

func (self *Handler) GetTransport() (trans Transport, err error) {
	if trans, err = self.tf.GetTransport(); err != nil {
		return
	}
	if trans, err = self.tw.Wraps(trans); err != nil {
		return
	}
	return
}
