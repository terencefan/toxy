package xtransport

import (
	"errors"
	"net"
)

const (
	sockTypeTcp  = "tcp"
	sockTypeUnix = "unix"
)

type TSocket struct {
	conn     net.Conn
	sockType string
}

type TSocketFactory struct {
	addr     string
	sockType string
}

func (self *TSocket) Read(message []byte) (int, error) {
	return self.conn.Read(message)
}

func (self *TSocket) Write(message []byte) (int, error) {
	return self.conn.Write(message)
}

func (self *TSocket) Close() error {
	return self.conn.Close()
}

func (self *TSocket) Flush() error {
	return nil
}

func (self *TSocketFactory) GetTransport() (Transport, error) {
	switch self.sockType {
	case sockTypeTcp:
		return NewTSocket(self.addr)
	case sockTypeUnix:
		return NewTUnixSocket(self.addr)
	default:
		return nil, errors.New("invalid socket type")
	}
}

func NewTSocket(addr string) (trans *TSocket, err error) {
	conn, err := net.Dial(sockTypeTcp, addr)
	if err != nil {
		return
	}
	trans = &TSocket{
		conn: conn,
	}
	return
}

func NewTUnixSocket(addr string) (trans *TSocket, err error) {
	conn, err := net.Dial(sockTypeUnix, addr)
	if err != nil {
		return
	}
	trans = &TSocket{
		conn: conn,
	}
	return
}

func NewTSocketConn(conn net.Conn) *TSocket {
	return &TSocket{
		conn: conn,
	}
}

func NewTSocketFactory(addr string) *TSocketFactory {
	return &TSocketFactory{
		addr:     addr,
		sockType: sockTypeTcp,
	}
}

func NewTUnixSocketFactory(addr string) *TSocketFactory {
	return &TSocketFactory{
		addr:     addr,
		sockType: sockTypeUnix,
	}
}
