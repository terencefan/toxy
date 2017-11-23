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

func (m *Messenger) ForwardMessage(iprot, oprot Protocol) (err error) {
	if err = m.ForwardMessageBegin(iprot, oprot); err != nil {
		return
	}
	if err = m.Forward(iprot, oprot, T_STRUCT); err != nil {
		return
	}
	if err = m.ForwardMessageEnd(iprot, oprot); err != nil {
		return
	}
	return
}

func (m *Messenger) Forward(iprot, oprot Protocol, ftype byte) error {
	if err := forward(iprot, oprot, ftype); err != nil {
		return err
	}
	return nil
}

func (m *Messenger) ForwardMessageBegin(iprot, oprot Protocol) error {
	name, mtype, seqid, err := iprot.ReadMessageBegin()
	if err != nil {
		return err
	}
	if err := oprot.WriteMessageBegin(name, mtype, seqid); err != nil {
		return err
	}
	return nil
}

func (m *Messenger) ForwardMessageEnd(iprot, oprot Protocol) error {
	if err := iprot.ReadMessageEnd(); err != nil {
		return err
	}
	if err := oprot.WriteMessageEnd(); err != nil {
		return err
	}
	if err := oprot.Flush(); err != nil {
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
	if err := m.iprot.Flush(); err != nil {
		return err
	}
	return
}

func (m *Messenger) FastReplyShutdown(name string, seqid int32) (err error) {

	ae := NewTApplicationException("toxy is shutting down", ExceptionUnknown)
	if err = m.skipIncomingMessage(); err != nil {
		return
	}
	if err = m.iprot.WriteMessageBegin(name, T_EXCEPTION, seqid); err != nil {
		return
	}
	if err = write_application_exception(ae, m.iprot); err != nil {
		return
	}
	if err = m.iprot.WriteMessageEnd(); err != nil {
		return
	}
	if err := m.iprot.Flush(); err != nil {
		return err
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
	name, mtype, seqid, err := m.iprot.ReadMessageBegin()
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
	iprot := NewStoredProtocol(m.iprot, name, T_CALL, seqid)
	if err := m.ForwardMessage(iprot, m.oprot); err != nil {
		panic(err)
	}
	if err := m.ForwardMessage(m.oprot, m.iprot); err != nil {
		panic(err)
	}
	return true
}
