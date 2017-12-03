package toxy

import (
	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/thrift"
)

type Messenger struct{}

func (m *Messenger) Reply(iprot, oprot Protocol) (err error) {
	if err = messenger.ForwardMessage(iprot, oprot); err != nil {
		return
	}
	if err = messenger.ForwardMessage(oprot, iprot); err != nil {
		return
	}
	return
}

func (m *Messenger) ForwardMessage(iprot, oprot Protocol) (err error) {
	if err = m.ForwardMessageBegin(iprot, oprot); err != nil {
		return
	}
	if err = forward(iprot, oprot, T_STRUCT); err != nil {
		return
	}
	if err = m.ForwardMessageEnd(iprot, oprot); err != nil {
		return
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
	if err := oprot.GetTransport().Flush(); err != nil {
		return err
	}
	return nil
}

func (m *Messenger) skipIncomingMessage(iprot Protocol) (err error) {
	if err = iprot.Skip(T_STRUCT); err != nil {
		return
	}
	if err = iprot.ReadMessageEnd(); err != nil {
		return
	}
	return
}

func (m *Messenger) FastReply(iprot Protocol, name string, seqid int32) (err error) {
	if err = m.skipIncomingMessage(iprot); err != nil {
		return
	}
	if err = iprot.WriteMessageBegin(name, T_REPLY, seqid); err != nil {
		return
	}
	if err = iprot.WriteByte(T_STOP); err != nil {
		return
	}
	if err = iprot.WriteMessageEnd(); err != nil {
		return
	}
	if err := iprot.GetTransport().Flush(); err != nil {
		return err
	}
	return
}

func (m *Messenger) FastReplyShutdown(iprot Protocol) (err error) {

	ae := NewTApplicationException(
		"toxy is shutting down",
		ExceptionUnknown,
	)
	if err = iprot.WriteMessageBegin("unknown", T_EXCEPTION, 0); err != nil {
		return
	}
	if err = WriteTApplicationException(iprot, ae); err != nil {
		return
	}
	if err = iprot.WriteMessageEnd(); err != nil {
		return
	}
	if err := iprot.GetTransport().Flush(); err != nil {
		return err
	}
	return
}

var messenger = &Messenger{}
