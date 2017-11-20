package xprotocol

// TODO use BigEndian

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	. "xthrift"
	. "xtransport"
)

const (
	VERSION_MASK uint32 = 0xffff0000
	VERSION_1    uint32 = 0x80010000
	TYPE_MASK    uint32 = 0x000000ff
)

var _ = T_STOP

const (
	maxMessageNameSize = 128
	maxBufSize         = 32
)

type TBinaryProtocol struct {
	trans       Transport
	strictRead  bool
	strictWrite bool
	buf         []byte
}

type TBinaryProtocolFactory struct {
	strictRead  bool
	strictWrite bool
}

func (self *TBinaryProtocolFactory) NewProtocol(t Transport) Protocol {
	return NewTBinaryProtocol(t, self.strictRead, self.strictWrite)
}

func NewTBinaryProtocol(conn Transport, sr, sw bool) *TBinaryProtocol {
	return &TBinaryProtocol{
		trans:       conn,
		strictRead:  sr,
		strictWrite: sw,
		buf:         make([]byte, maxBufSize),
	}
}

func NewTBinaryProtocolFactory(strictRead, strictWrite bool) *TBinaryProtocolFactory {
	return &TBinaryProtocolFactory{
		strictRead:  strictRead,
		strictWrite: strictWrite,
	}
}

func (p *TBinaryProtocol) getBuf(ln int32) []byte {
	b := p.buf
	if ln > int32(len(b)) {
		b = make([]byte, ln)
	} else {
		b = b[:ln]
	}
	return b
}

func (p *TBinaryProtocol) WriteMessageBegin(name string, messageType byte, seqid int32) error {
	if p.strictWrite {
		if err := p.WriteI32(int32(VERSION_1 | uint32(messageType))); err != nil {
			return err
		}
		if err := p.WriteString(name); err != nil {
			return err
		}
	} else {
		if err := p.WriteString(name); err != nil {
			return err
		}
		if err := p.WriteByte(messageType); err != nil {
			return err
		}
	}
	return p.WriteI32(seqid)
}

func (p *TBinaryProtocol) WriteMessageEnd() error {
	return nil
}

func (p *TBinaryProtocol) WriteStructBegin(name string) error {
	return nil
}

func (p *TBinaryProtocol) WriteStructEnd() error {
	return nil
}

func (p *TBinaryProtocol) WriteFieldBegin(name string, fieldType byte, id int16) error {
	if err := p.WriteByte(fieldType); err != nil {
		return err
	}
	return p.WriteI16(id)
}

func (p *TBinaryProtocol) WriteFieldEnd() error {
	return nil
}

func (p *TBinaryProtocol) WriteFieldStop() error {
	return p.WriteByte(T_STOP)
}

func (p *TBinaryProtocol) WriteMapBegin(ktype byte, vtype byte, size int) error {
	if err := p.WriteByte(ktype); err != nil {
		return err
	}
	if err := p.WriteByte(vtype); err != nil {
		return err
	}
	return p.WriteI32(int32(size))
}

func (p *TBinaryProtocol) WriteMapEnd() error {
	return nil
}

func (p *TBinaryProtocol) WriteListBegin(etype byte, size int) error {
	if err := p.WriteByte(etype); err != nil {
		return err
	}
	return p.WriteI32(int32(size))
}

func (p *TBinaryProtocol) WriteListEnd() error {
	return nil
}

func (p *TBinaryProtocol) WriteSetBegin(etype byte, size int) error {
	if err := p.WriteByte(etype); err != nil {
		return err
	}
	return p.WriteI32(int32(size))
}

func (p *TBinaryProtocol) WriteSetEnd() error {
	return nil
}

func (p *TBinaryProtocol) WriteBool(value bool) error {
	if value {
		return p.WriteByte(1)
	} else {
		return p.WriteByte(0)
	}
}

func (p *TBinaryProtocol) WriteByte(value byte) error {
	b := p.buf
	if b == nil {
		b = []byte{value}
	} else {
		b[0] = value
	}
	_, err := p.trans.Write(b[:1])
	return err
}

func (p *TBinaryProtocol) WriteI16(value int16) error {
	b := p.buf
	if b == nil {
		b = []byte{0, 0}
	}
	binary.BigEndian.PutUint16(b, uint16(value))
	_, err := p.trans.Write(b[:2])
	return err
}

func (p *TBinaryProtocol) WriteI32(value int32) error {
	b := p.buf
	if b == nil {
		b = []byte{0, 0, 0, 0}
	}
	binary.BigEndian.PutUint32(b, uint32(value))
	_, err := p.trans.Write(b[:4])
	return err
}

func (p *TBinaryProtocol) WriteI64(value int64) error {
	b := p.buf
	if b == nil {
		b = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	}
	binary.BigEndian.PutUint64(b, uint64(value))
	_, err := p.trans.Write(b[:8])
	return err
}

func (p *TBinaryProtocol) WriteDouble(value float64) error {
	b := p.buf
	if b == nil {
		b = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	}
	binary.BigEndian.PutUint64(b, math.Float64bits(value))
	_, err := p.trans.Write(b[:8])
	return err
}

func (p *TBinaryProtocol) WriteString(value string) error {
	ln := int32(len(value))

	b := p.getBuf(ln)
	if err := p.WriteI32(ln); err != nil {
		return err
	}
	copy(b, value)
	_, err := p.trans.Write(b)
	return err
}

func (p *TBinaryProtocol) WriteBytes(value []byte) error {
	if err := p.WriteI32(int32(len(value))); err != nil {
		return err
	}
	_, err := p.trans.Write(value)
	return err
}

func (p *TBinaryProtocol) ReadMessageBegin() (
	name string, messageType byte, seqid int32, err error,
) {
	size, err := p.ReadI32()
	if err != nil {
		return
	}

	if size < 0 {
		version := uint32(size) & VERSION_MASK
		if version != VERSION_1 {
			err = ProtocolError{
				"BinaryProtocol",
				"bad version in ReadMessageBegin",
			}
			return
		}
		messageType = byte(uint32(size) & TYPE_MASK)
		if name, err = p.ReadString(); err != nil {
			return
		}
	} else {
		if p.strictRead {
			err = ProtocolError{
				"BinaryProtocol",
				"no protocol version header",
			}
			return
		}
		if size > maxMessageNameSize {
			err = ProtocolError{
				"BinaryProtocol",
				"message name exceeds max size",
			}
			return
		}
		nameBytes := make([]byte, size)
		if _, err = p.trans.Read(nameBytes); err != nil {
			return
		}
		name = string(nameBytes)
		if messageType, err = p.ReadByte(); err != nil {
			return
		}
	}
	seqid, err = p.ReadI32()
	return
}

func (p *TBinaryProtocol) ReadMessageEnd() error {
	return nil
}

func (p *TBinaryProtocol) ReadStructBegin() error {
	return nil
}

func (p *TBinaryProtocol) ReadStructEnd() error {
	return nil
}

func (p *TBinaryProtocol) ReadFieldBegin() (
	ftype byte, id int16, err error,
) {
	if ftype, err = p.ReadByte(); err != nil {
		return
	} else if ftype == T_STOP {
		return
	}
	id, err = p.ReadI16()
	return
}

func (p *TBinaryProtocol) ReadFieldEnd() error {
	return nil
}

func (p *TBinaryProtocol) ReadMapBegin() (
	ktype byte, vtype byte, size int, err error,
) {
	if ktype, err = p.ReadByte(); err != nil {
		return
	}
	if vtype, err = p.ReadByte(); err != nil {
		return
	}
	sz, err := p.ReadI32()
	size = int(sz)
	return
}

func (p *TBinaryProtocol) ReadMapEnd() error {
	return nil
}

func (p *TBinaryProtocol) ReadListBegin() (
	etype byte, size int, err error,
) {
	if etype, err = p.ReadByte(); err != nil {
		return
	}
	sz, err := p.ReadI32()
	size = int(sz)
	return
}

func (p *TBinaryProtocol) ReadListEnd() error {
	return nil
}

func (p *TBinaryProtocol) ReadSetBegin() (
	etype byte, size int, err error,
) {
	return p.ReadListBegin()
}

func (p *TBinaryProtocol) ReadSetEnd() error {
	return nil
}

func (p *TBinaryProtocol) ReadBool() (bool, error) {
	if b, err := p.ReadByte(); err != nil {
		return false, err
	} else if b != 0 {
		return true, nil
	}
	return false, nil
}

func (p *TBinaryProtocol) ReadByte() (value byte, err error) {
	_, err = io.ReadFull(p.trans, p.buf[:1])
	value = p.buf[0]
	return
}

func (p *TBinaryProtocol) ReadI16() (value int16, err error) {
	_, err = io.ReadFull(p.trans, p.buf[:2])
	value = int16(binary.BigEndian.Uint16(p.buf))
	return
}

func (p *TBinaryProtocol) ReadI32() (value int32, err error) {
	_, err = io.ReadFull(p.trans, p.buf[:4])
	value = int32(binary.BigEndian.Uint32(p.buf))
	return
}

func (p *TBinaryProtocol) ReadI64() (value int64, err error) {
	_, err = io.ReadFull(p.trans, p.buf[:8])
	value = int64(binary.BigEndian.Uint64(p.buf))
	return
}

func (p *TBinaryProtocol) ReadDouble() (value float64, err error) {
	_, err = io.ReadFull(p.trans, p.buf[:8])
	value = math.Float64frombits(binary.BigEndian.Uint64(p.buf))
	return
}

func (p *TBinaryProtocol) ReadString() (s string, err error) {
	ln, err := p.ReadI32()
	if err != nil || ln == 0 {
		return
	}

	b, err := p.ReadBytes(ln)
	if err != nil {
		return
	}
	return string(b), nil
}

func (p *TBinaryProtocol) ReadBytes(ln int32) ([]byte, error) {
	if ln == 0 {
		return nil, nil
	} else if ln < 0 {
		return nil, ProtocolError{
			"BinaryProtocol",
			"negative length while reading bytes",
		}
	}

	b := p.getBuf(ln)
	if _, err := io.ReadFull(p.trans, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (p *TBinaryProtocol) skip(ln int32) (err error) {
	if ln <= 0 {
		return
	}
	b := p.getBuf(ln)
	if _, err = io.ReadFull(p.trans, b); err != nil {
		return
	}
	return
}

func (p *TBinaryProtocol) skipString() error {
	if ln, err := p.ReadI32(); err != nil {
		return err
	} else {
		return p.skip(ln)
	}
}

func (p *TBinaryProtocol) skipList() error {
	etype, size, err := p.ReadListBegin()
	if err != nil {
		return err
	}
	for i := 0; i < size; i++ {
		if err := p.Skip(etype); err != nil {
			return err
		}
	}
	return nil
}

func (p *TBinaryProtocol) skipMap() error {
	ktype, vtype, size, err := p.ReadMapBegin()
	if err != nil {
		return err
	}
	for i := 0; i < size; i++ {
		if err := p.Skip(ktype); err != nil {
			return err
		}
		if err := p.Skip(vtype); err != nil {
			return err
		}
	}
	return nil
}

func (p *TBinaryProtocol) skipStruct() error {
	for {
		ftype, _, err := p.ReadFieldBegin()
		if err != nil {
			return err
		} else if ftype == T_STOP {
			return nil
		}
		if err := p.Skip(ftype); err != nil {
			return err
		}
	}
}

func (p *TBinaryProtocol) Skip(ftype byte) error {
	switch ftype {
	case T_BOOL:
		return p.skip(1)
	case T_I08:
		return p.skip(1)
	case T_I16:
		return p.skip(2)
	case T_I32:
		return p.skip(4)
	case T_I64:
		return p.skip(8)
	case T_DOUBLE:
		return p.skip(8)
	case T_STRING:
		return p.skipString()
	case T_SET:
		return p.skipList()
	case T_LIST:
		return p.skipList()
	case T_MAP:
		return p.skipMap()
	case T_STRUCT:
		return p.skipStruct()
	default:
		return errors.New("unsupported field type")
	}
}

func (p *TBinaryProtocol) Flush() error {
	return p.trans.Flush()
}

func (p *TBinaryProtocol) Close() error {
	return p.trans.Close()
}
