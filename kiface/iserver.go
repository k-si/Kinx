package kiface

// 服务器接口
type IServer interface {
	Start()
	Serve() error
	Stop()
	Recycle()
	GetMsgHandler() IMsgHandler
	GetConnMgr() IConnectionManager
	SetAfterConnSuccess(func(IConnection)) IServer
	SetBeforeConnDestroy(func(IConnection)) IServer
	CallAfterConnSuccess(IConnection)
	CallBeforeConnDestroy(IConnection)
	AddRouter(uint32, IRouter) IServer
}
