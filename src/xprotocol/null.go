package xprotocol

type t_null_protocol struct {
}

var TNullProtocol = &t_null_protocol{}

func (p *t_null_protocol) WriteMessageBegin(name string, messageType byte, seqid int32) error {
	return nil
}

func (p *t_null_protocol) WriteMessageEnd() error {
	return nil
}

func (p *t_null_protocol) WriteStructBegin(name string) error {
	return nil
}

func (p *t_null_protocol) WriteStructEnd() error {
	return nil
}

func (p *t_null_protocol) WriteFieldBegin(name string, fieldType byte, id int16) error {
	return nil
}

func (p *t_null_protocol) WriteFieldEnd() error {
	return nil
}

func (p *t_null_protocol) WriteFieldStop() error {
	return nil
}

func (p *t_null_protocol) WriteMapBegin(ktype byte, vtype byte, size int) error {
	return nil
}

func (p *t_null_protocol) WriteMapEnd() error {
	return nil
}

func (p *t_null_protocol) WriteListBegin(etype byte, size int) error {
	return nil
}

func (p *t_null_protocol) WriteListEnd() error {
	return nil
}

func (p *t_null_protocol) WriteSetBegin(etype byte, size int) error {
	return nil
}

func (p *t_null_protocol) WriteSetEnd() error {
	return nil
}

func (p *t_null_protocol) WriteBool(value bool) error {
	return nil
}

func (p *t_null_protocol) WriteByte(value byte) error {
	return nil
}

func (p *t_null_protocol) WriteI16(value int16) error {
	return nil
}

func (p *t_null_protocol) WriteI32(value int32) error {
	return nil
}

func (p *t_null_protocol) WriteI64(value int64) error {
	return nil
}

func (p *t_null_protocol) WriteDouble(value float64) error {
	return nil
}

func (p *t_null_protocol) WriteString(value string) error {
	return nil
}

func (p *t_null_protocol) ReadMessageBegin() (
	name string, messageType byte, seqid int32, err error,
) {
	return
}

func (p *t_null_protocol) ReadMessageEnd() error {
	return nil
}

func (p *t_null_protocol) ReadStructBegin() error {
	return nil
}

func (p *t_null_protocol) ReadStructEnd() error {
	return nil
}

func (p *t_null_protocol) ReadFieldBegin() (
	ftype byte, id int16, err error,
) {
	return
}

func (p *t_null_protocol) ReadFieldEnd() error {
	return nil
}

func (p *t_null_protocol) ReadMapBegin() (
	ktype byte, vtype byte, size int, err error,
) {
	return
}

func (p *t_null_protocol) ReadMapEnd() error {
	return nil
}

func (p *t_null_protocol) ReadListBegin() (
	etype byte, size int, err error,
) {
	return
}

func (p *t_null_protocol) ReadListEnd() error {
	return nil
}

func (p *t_null_protocol) ReadSetBegin() (
	etype byte, size int, err error,
) {
	return
}

func (p *t_null_protocol) ReadSetEnd() error {
	return nil
}

func (p *t_null_protocol) ReadBool() (bool, error) {
	return false, nil
}

func (p *t_null_protocol) ReadByte() (value byte, err error) {
	return
}

func (p *t_null_protocol) ReadI16() (value int16, err error) {
	return
}

func (p *t_null_protocol) ReadI32() (value int32, err error) {
	return
}

func (p *t_null_protocol) ReadI64() (value int64, err error) {
	return
}

func (p *t_null_protocol) ReadDouble() (value float64, err error) {
	return
}

func (p *t_null_protocol) ReadString() (s string, err error) {
	return
}

func (p *t_null_protocol) Skip(ftype byte) error {
	return nil
}

func (p *t_null_protocol) Flush() error {
	return nil
}

func (p *t_null_protocol) Close() error {
	return nil
}
