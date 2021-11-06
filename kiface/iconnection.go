package kiface

import "net"

type Iconnection interface {
	Start()
	Stop()
}

type HandleFunc func(*net.TCPConn, []byte, int) error
