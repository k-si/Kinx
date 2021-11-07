package knet

import (
	"kinx/kiface"
)

type Request struct {
	conn kiface.IConnection
	data []byte
}

func (r *Request) GetConnection() kiface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}

func NewRequest(conn kiface.IConnection, data []byte) kiface.IRequest{
	req := &Request{
		conn: conn,
		data: data,
	}
	return req
}