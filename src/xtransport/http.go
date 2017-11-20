package xtransport

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

type THttpTransport struct {
	addr string
	path string
	buf  *bytes.Buffer
}

type THttpTransportFactory struct {
	addr string
	path string
}

// TODO write a real http wrapper
type THttpTransportWrapper struct {
	path string
}

func (self *THttpTransport) Read(message []byte) (int, error) {
	return self.buf.Read(message)
}

func (self *THttpTransport) Write(message []byte) (int, error) {
	return self.buf.Write(message)
}

func (self *THttpTransport) Close() error {
	return nil
}

func (self *THttpTransport) Flush() (err error) {
	uri := fmt.Sprintf("http://%s%s", self.addr, self.path)

	resp, err := http.Post(uri, "application/thrift", self.buf)

	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	if _, err := self.buf.ReadFrom(resp.Body); err != nil {
		return err
	}
	return nil
}

func (self *THttpTransportFactory) GetTransport() (
	trans Transport, err error,
) {
	trans = &THttpTransport{
		addr: self.addr,
		path: self.path,
		buf:  bytes.NewBuffer([]byte{}),
	}
	return
}

func NewTHttpTransportFactory(addr, path string) *THttpTransportFactory {
	return &THttpTransportFactory{
		addr: addr,
		path: path,
	}
}
