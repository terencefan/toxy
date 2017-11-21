package xtransport

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
	. "xception"
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

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Post(uri, "application/thrift", self.buf)

	if ne, ok := err.(net.Error); ok && ne.Timeout() {
		err = NewTApplicationException("toxy timeout exceeded", 0)
	}
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
