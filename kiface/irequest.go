package kiface

type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
}