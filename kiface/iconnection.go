package kiface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	SendMessage(uint32, []byte) error
	GetConnectionID() uint32
	SetProperty(string, interface{})
	GetProperty(string) (interface{}, error)
	RemoveProperty(string) error
}
