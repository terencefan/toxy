package xtransport

import (
	"bufio"
)

const (
	DEFAULT_BUF_SIZE = 4096
)

type TBufferedTransport struct {
	trans Transport
	rbuf  *bufio.Reader
	wbuf  *bufio.Writer
}

type TBufferedTransportFactory struct {
	rbufsize int
	wbufsize int
}

func (t *TBufferedTransport) Read(message []byte) (int, error) {
	return t.rbuf.Read(message)
}

func (t *TBufferedTransport) Write(message []byte) (int, error) {
	return t.wbuf.Write(message)
}

func (t *TBufferedTransport) Flush() (err error) {
	if err = t.wbuf.Flush(); err != nil {
		return
	}
	if err = t.trans.Flush(); err != nil {
		return
	}
	return
}

func (t *TBufferedTransport) Close() (err error) {
	t.rbuf = nil
	t.wbuf = nil
	if err = t.trans.Close(); err != nil {
		return
	}
	t.trans = nil
	return
}

func NewTBufferedTransport(trans Transport) *TBufferedTransport {
	return &TBufferedTransport{
		trans: trans,
		rbuf:  bufio.NewReaderSize(trans, DEFAULT_BUF_SIZE),
		wbuf:  bufio.NewWriterSize(trans, DEFAULT_BUF_SIZE),
	}
}

func NewTBufferedTransportSize(trans Transport, rbufsize, wbufsize int) *TBufferedTransport {
	return &TBufferedTransport{
		trans: trans,
		rbuf:  bufio.NewReaderSize(trans, rbufsize),
		wbuf:  bufio.NewWriterSize(trans, wbufsize),
	}
}

func (t *TBufferedTransportFactory) Wraps(trans Transport) (Transport, error) {
	return NewTBufferedTransportSize(trans, t.rbufsize, t.wbufsize), nil
}

func NewTBufferedTransportFactory(rbufsize, wbufsize int) *TBufferedTransportFactory {
	return &TBufferedTransportFactory{
		rbufsize: rbufsize,
		wbufsize: wbufsize,
	}
}
