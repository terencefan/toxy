package toxy

import (
	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/thrift"
)

type Messenger struct{}

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
	if err := oprot.Flush(); err != nil {
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
	if err := iprot.Flush(); err != nil {
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
	if err = write_application_exception(ae, iprot); err != nil {
		return
	}
	if err = iprot.WriteMessageEnd(); err != nil {
		return
	}
	if err := iprot.Flush(); err != nil {
		return err
	}
	return
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

var messenger = &Messenger{}

func read_message_begin(iprot Protocol) (name string, seqid int32) {
	name, mtype, seqid, err := iprot.ReadMessageBegin()
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

func fast_reply(iprot Protocol, name string, seqid int32) bool {
	if err := messenger.FastReply(iprot, name, seqid); err != nil {
		panic(err)
	}
	return true
}

func fast_reply_shutdown(iprot Protocol) bool {
	if err := messenger.FastReplyShutdown(iprot); err != nil {
		panic(err)
	}
	return false
}

func reply(iprot, oprot Protocol) bool {
	if err := messenger.ForwardMessage(iprot, oprot); err != nil {
		panic(err)
	}
	if err := messenger.ForwardMessage(oprot, iprot); err != nil {
		panic(err)
	}
	return true
}
