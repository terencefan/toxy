package xprotocol

type StoredProtocol struct {
	Protocol
	name  string
	mtype byte
	seqid int32
}

func NewStoredProtocol(p Protocol, name string, mtype byte, seqid int32) *StoredProtocol {
	return &StoredProtocol{
		Protocol: p,
		name:     name,
		mtype:    mtype,
		seqid:    seqid,
	}
}

func (p StoredProtocol) ReadMessageBegin() (string, byte, int32, error) {
	return p.name, p.mtype, p.seqid, nil
}
