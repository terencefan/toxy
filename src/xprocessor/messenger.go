package xprocessor

import (
	. "xception"
	"xlog"
	. "xprotocol"
	. "xthrift"
)

type Messenger struct {
	iprot Protocol
	oprot Protocol
}

type MessengerProtocolWriter interface {
	Write(Protocol) error
}

func (m *Messenger) Reverse() {
	m.iprot, m.oprot = m.oprot, m.iprot
}

func (m *Messenger) ReadMessageBegin() (name string, mtype byte, seqid int32, err error) {
	return m.iprot.ReadMessageBegin()
}

func (m *Messenger) WriteMessageBegin(name string, mtype byte, seqid int32) error {
	return m.oprot.WriteMessageBegin(name, mtype, seqid)
}

func (m *Messenger) Forward(ftype byte) error {
	if err := forward(m.iprot, m.oprot, ftype); err != nil {
		return err
	}
	return nil
}

func (m *Messenger) ForwardMessageBegin() error {
	name, mtype, seqid, err := m.iprot.ReadMessageBegin()
	if err != nil {
		return err
	}
	if err := m.oprot.WriteMessageBegin(name, mtype, seqid); err != nil {
		return err
	}
	return nil
}

func (m *Messenger) ForwardMessageEnd() error {
	if err := m.iprot.ReadMessageEnd(); err != nil {
		return err
	}
	if err := m.oprot.WriteMessageEnd(); err != nil {
		return err
	}
	if err := m.oprot.Flush(); err != nil {
		return err
	}
	return nil
}

func (m *Messenger) skipIncomingMessage() (err error) {
	if err = m.iprot.Skip(T_STRUCT); err != nil {
		return
	}
	if err = m.iprot.ReadMessageEnd(); err != nil {
		return
	}
	return
}

func (m *Messenger) writeApplicationException(name string, seqid int32, ae *TApplicationException) (err error) {
	if err = m.iprot.WriteMessageBegin(name, T_EXCEPTION, seqid); err != nil {
		return
	}
	if err = write_application_exception(ae, m.iprot); err != nil {
		return
	}
	if err = m.iprot.WriteMessageEnd(); err != nil {
		return
	}
	return
}

func (m *Messenger) FastReply(seqid int32) (err error) {
	xlog.Debug("fast reply: ping")
	if err = m.skipIncomingMessage(); err != nil {
		return
	}
	if err = m.iprot.WriteMessageBegin("ping", T_REPLY, seqid); err != nil {
		return
	}
	if err = m.iprot.WriteByte(T_STOP); err != nil {
		return
	}
	if err = m.iprot.WriteMessageEnd(); err != nil {
		return
	}
	return
}

func (m *Messenger) FastReplyShutdown(name string, seqid int32) (err error) {
	ae := NewTApplicationException("toxy is shutting down", ExceptionUnknown)
	if err = m.skipIncomingMessage(); err != nil {
		return
	}
	if err = m.writeApplicationException(name, seqid, ae); err != nil {
		return
	}
	return
}

func (m *Messenger) SetOutputProtocol(oprot Protocol) {
	m.oprot = oprot
}

func (m *Messenger) DelOutputProtocol() {
	m.oprot.Close()
	m.oprot = nil
}

func NewMessenger(protocol Protocol) *Messenger {
	return &Messenger{
		iprot: protocol,
	}
}

func write_application_exception(e *TApplicationException, proto Protocol) (err error) {
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

func read_header(m *Messenger) (name string, seqid int32) {
	xlog.Debug("read message header")
	name, mtype, seqid, err := m.ReadMessageBegin()
	if err != nil {
		panic(err)
	} else if mtype == T_ONEWAY {
		// TODO reply exception "doesn't support oneway request yet."
	} else if mtype != T_CALL {
		// TODO raise exception.
	} else {
		return
	}
	return
}

func fast_reply(m *Messenger, seqid int32) bool {
	if err := m.FastReply(seqid); err != nil {
		panic(err)
	}
	return true
}

func fast_reply_shutdown(m *Messenger, name string, seqid int32) bool {
	if err := m.FastReplyShutdown(name, seqid); err != nil {
		panic(err)
	}
	return false
}

func reply(m *Messenger, name string, seqid int32) bool {
	if err := m.WriteMessageBegin(name, T_CALL, seqid); err != nil {
		panic(err)
	}
	if err := m.Forward(T_STRUCT); err != nil {
		panic(err)
	}
	if err := m.ForwardMessageEnd(); err != nil {
		panic(err)
	}
	m.Reverse()

	if err := m.ForwardMessageBegin(); err != nil {
		panic(err)
	}
	if err := m.Forward(T_STRUCT); err != nil {
		panic(err)
	}
	if err := m.ForwardMessageEnd(); err != nil {
		panic(err)
	}
	m.Reverse()
	return true
}
