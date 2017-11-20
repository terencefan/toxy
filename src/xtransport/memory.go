package xtransport

import "bytes"

type TMemoryBuffer struct {
	*bytes.Buffer
}

func (self *TMemoryBuffer) Close() error {
	return nil
}

func (self *TMemoryBuffer) Flush() error {
	return nil
}

func (self *TMemoryBuffer) GetBytes() []byte {
	return self.Buffer.Bytes()
}

func NewTMemoryBuffer() *TMemoryBuffer {
	return &TMemoryBuffer{
		bytes.NewBuffer([]byte{}),
	}
}
