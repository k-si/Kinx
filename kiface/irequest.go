package kiface

type IRequest interface {
	GetConnection() IConnection
	GetMsg() IMessage
}
