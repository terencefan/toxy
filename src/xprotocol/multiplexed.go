package xprotocol

import "fmt"

type TMultiplexedProtocol struct {
	Protocol
	service   string
	delimeter string
}

func NewTMultiplexedProtocol(p Protocol, service string) *TMultiplexedProtocol {
	return &TMultiplexedProtocol{
		Protocol:  p,
		service:   service,
		delimeter: ":",
	}
}

func (p TMultiplexedProtocol) WriteMessageBegin(
	name string, mtype byte, seqid int32,
) error {
	name = fmt.Sprintf("%s%s%s", p.service, p.delimeter, name)
	return p.Protocol.WriteMessageBegin(name, mtype, seqid)
}
