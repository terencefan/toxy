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

func new_http_tf(section *ini.Section) (tf TransportFactory, err error) {
	// TODO check if addr / path exist in config
	addr := section.Key("addr").String()
	path := section.Key("path").String()
	tf = NewTHttpTransportFactory(addr, path)
	return
}

func new_socket_tf(section *ini.Section) (tf TransportFactory, err error) {
	return
}

func new_buffered_tw(section *ini.Section) (tw TransportWrapper, err error) {
	rbufsize := 4096
	wbufsize := 4096
	tw = NewTBufferedTransportFactory(rbufsize, wbufsize)
	return
}

func new_framed_tw(section *ini.Section) (tw TransportWrapper, err error) {
	err = errors.New("framed transport wrapper has not been supported yet")
	return
}

func NewHandler(section *ini.Section) (h *Handler, err error) {
	h = &Handler{
		tw: TTransportWrapper,
	}

	// transport
	if section.HasKey("transport") {
		transport := section.Key("transport").String()
		switch transport {
		case "http":
			h.tf, err = new_http_tf(section)
		case "socket":
			h.tf, err = new_socket_tf(section)
		default:
			// TODO more error messages.
			err = errors.New("invalid transport: " + transport)
			return nil, err
		}
	} else {
		err = errors.New("transport has not been defined")
	}

	if err != nil {
		return nil, err
	}

	// wrapper
	if section.HasKey("wrapper") {
		wrapper := section.Key("wrapper").String()
		switch wrapper {
		case "buffered":
			h.tw, err = new_buffered_tw(section)
		case "framed":
			h.tw, err = new_framed_tw(section)
		default:
			err = errors.New("invalid transport wrapper: " + wrapper)
		}
	} else {
		h.tw = TTransportWrapper
	}

	if err != nil {
		return nil, err
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
