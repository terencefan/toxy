package xprocessor

import (
	"errors"
	. "xprotocol"
	. "xthrift"
)

func forward(iprot, oprot Protocol, ftype byte) error {
	switch ftype {
	case T_BOOL:
		return forward_bool(iprot, oprot)
	case T_I08:
		return forward_byte(iprot, oprot)
	case T_I16:
		return forward_i16(iprot, oprot)
	case T_I32:
		return forward_i32(iprot, oprot)
	case T_I64:
		return forward_i64(iprot, oprot)
	case T_DOUBLE:
		return forward_double(iprot, oprot)
	case T_STRING:
		return forward_string(iprot, oprot)
	case T_SET:
		return forward_set(iprot, oprot)
	case T_LIST:
		return forward_list(iprot, oprot)
	case T_MAP:
		return forward_map(iprot, oprot)
	case T_STRUCT:
		return forward_struct(iprot, oprot)
	default:
		return errors.New("unsupported field type")
	}
}

func forward_bool(iprot, oprot Protocol) error {
	b, err := iprot.ReadBool()
	if err != nil {
		return err
	}
	return oprot.WriteBool(b)
}

func forward_byte(iprot, oprot Protocol) error {
	b, err := iprot.ReadByte()
	if err != nil {
		return err
	}
	return oprot.WriteByte(b)
}

func forward_i16(iprot, oprot Protocol) error {
	i, err := iprot.ReadI16()
	if err != nil {
		return err
	}
	return oprot.WriteI16(i)
}

func forward_i32(iprot, oprot Protocol) error {
	i, err := iprot.ReadI32()
	if err != nil {
		return err
	}
	return oprot.WriteI32(i)
}

func forward_i64(iprot, oprot Protocol) error {
	i, err := iprot.ReadI64()
	if err != nil {
		return err
	}
	return oprot.WriteI64(i)
}

func forward_double(iprot, oprot Protocol) error {
	d, err := iprot.ReadDouble()
	if err != nil {
		return err
	}
	return oprot.WriteDouble(d)
}

func forward_string(iprot, oprot Protocol) error {
	s, err := iprot.ReadString()
	if err != nil {
		return err
	}
	return oprot.WriteString(s)
}

func forward_set(iprot, oprot Protocol) error {
	etype, size, err := iprot.ReadSetBegin()
	if err != nil {
		return err
	}
	if err := oprot.WriteSetBegin(etype, size); err != nil {
		return err
	}
	for i := 0; i < size; i++ {
		if err := forward(iprot, oprot, etype); err != nil {
			return err
		}
	}
	if err := iprot.ReadSetEnd(); err != nil {
		return err
	}
	return oprot.WriteSetEnd()
}

func forward_list(iprot, oprot Protocol) error {
	etype, size, err := iprot.ReadListBegin()
	if err != nil {
		return err
	}
	if err := oprot.WriteListBegin(etype, size); err != nil {
		return err
	}
	for i := 0; i < size; i++ {
		if err := forward(iprot, oprot, etype); err != nil {
			return err
		}
	}
	if err := iprot.ReadListEnd(); err != nil {
		return err
	}
	return oprot.WriteListEnd()
}

func forward_map(iprot, oprot Protocol) error {
	ktype, vtype, size, err := iprot.ReadMapBegin()
	if err != nil {
		return err
	}
	if err := oprot.WriteMapBegin(ktype, vtype, size); err != nil {
		return err
	}
	for i := 0; i < size; i++ {
		if err := forward(iprot, oprot, ktype); err != nil {
			return err
		}
		if err := forward(iprot, oprot, vtype); err != nil {
			return err
		}
	}
	if err := iprot.ReadMapEnd(); err != nil {
		return err
	}
	return oprot.WriteMapEnd()
}

func forward_struct(iprot, oprot Protocol) error {
	if err := iprot.ReadStructBegin(); err != nil {
		return err
	}
	if err := oprot.WriteStructBegin("hi"); err != nil {
		return err
	}
	if err := forward_fields(iprot, oprot); err != nil {
		return err
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return err
	}
	return oprot.WriteStructEnd()
}

func forward_fields(iprot, oprot Protocol) error {
	for {
		ftype, fid, err := iprot.ReadFieldBegin()
		if err != nil {
			return err
		} else if ftype == T_STOP {
			if err := oprot.WriteFieldStop(); err != nil {
				return err
			}
			return nil
		} else {
			if err := oprot.WriteFieldBegin("hi", ftype, fid); err != nil {
				return err
			}
		}
		forward(iprot, oprot, ftype)
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return err
		}
	}
}
