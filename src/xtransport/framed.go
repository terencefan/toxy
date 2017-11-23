package xtransport

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	MAX_FRAME_SIZE = 1024 * 1024
)

type TFramedTransport struct {
	trans   Transport
	rframed bool
	wframed bool
	rbuf    *bytes.Buffer
	wbuf    *bytes.Buffer
}

type TFramedTransportFactory struct {
	rframed bool
	wframed bool
}

type ErrFrameInvalid struct {
}

func (e ErrFrameInvalid) Error() string {
	return "frame size invalid!"
}

type ErrFrameTooBig struct {
	size int64
}

func (e ErrFrameTooBig) Error() string {
	return fmt.Sprintf(
		"frame size exceeded: (%d > %d)",
		e.size,
		MAX_FRAME_SIZE,
	)
}

func (t *TFramedTransport) readFrame() (err error) {
	return errors.New("Not implement error")
}

func (t *TFramedTransport) Read(message []byte) (int, error) {
	if t.rframed {
		if t.rbuf.Len() == 0 {
			if err := t.readFrame(); err != nil {
				return 0, err
			}
		}
		return t.rbuf.Read(message)
	} else {
		return t.trans.Read(message)
	}
}

func (t *TFramedTransport) Write(message []byte) (int, error) {
	if t.wframed {
		l, err := t.wbuf.Write(message)
		if err != nil {
			return 0, err
		}
		if t.wbuf.Len() > MAX_FRAME_SIZE {
			return 0, &ErrFrameTooBig{int64(t.wbuf.Len())}
		}
		return l, nil
	} else {
		return t.trans.Write(message)
	}
}

func (t *TFramedTransport) Flush() (err error) {
	if t.wframed {
		size := uint32(t.wbuf.Len())
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, size)
		if _, err = t.trans.Write(buf); err != nil {
			return
		}
		if _, err = t.wbuf.WriteTo(t.trans); err != nil {
			return
		}
	}
	return t.trans.Flush()
}

func (t *TFramedTransport) Close() (err error) {
	t.rbuf = nil
	t.wbuf = nil
	if err = t.trans.Close(); err != nil {
		return
	}
	t.trans = nil
	return
}

func NewTFramedTransport(trans Transport, rframed, wframed bool) *TFramedTransport {
	return &TFramedTransport{
		trans:   trans,
		rframed: rframed,
		wframed: wframed,
		rbuf:    new(bytes.Buffer),
		wbuf:    new(bytes.Buffer),
	}
}

func (t *TFramedTransportFactory) Wraps(trans Transport) (Transport, error) {
	return NewTFramedTransport(trans, t.rframed, t.wframed), nil
}

func NewTFramedTransportFactory(rframed, wframed bool) *TFramedTransportFactory {
	return &TFramedTransportFactory{
		rframed: rframed,
		wframed: wframed,
	}
}
