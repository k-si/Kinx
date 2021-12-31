package knet

import (
	"github.com/k-si/Kinx/kiface"
)

type Request struct {
	conn kiface.IConnection
	msg kiface.IMessage
}

func (r *Request) GetConnection() kiface.IConnection {
	return r.conn
}

func (r *Request) GetMsg() kiface.IMessage {
	return r.msg
}

func NewRequest(conn kiface.IConnection, msg kiface.IMessage) kiface.IRequest{
	req := &Request{
		conn: conn,
		msg: msg,
	}
	return req
}