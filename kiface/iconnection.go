package kiface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	SendMessage(uint32, []byte) error
}
