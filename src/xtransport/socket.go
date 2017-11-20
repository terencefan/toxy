package xtransport

import (
	"fmt"
	"net"
)

type TSocket struct {
	conn net.Conn
}

type TSocketFactory struct {
	addr string
}

func (self *TSocket) Read(message []byte) (int, error) {
	return self.conn.Read(message)
}

func (self *TSocket) Write(message []byte) (int, error) {
	fmt.Println(message)
	return self.conn.Write(message)
}

func (self *TSocket) Close() error {
	return self.conn.Close()
}

func (self *TSocket) Flush() error {
	return nil
}

func (self *TSocketFactory) GetTransport() (Transport, error) {
	return NewTSocket(self.addr)
}

func NewTSocket(addr string) (trans *TSocket, err error) {
	conn, err := net.Dial("tcp", addr)
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
		addr: addr,
	}
}
