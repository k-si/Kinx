package kiface

// 服务器接口
type IServer interface {
	Start()
	Serve()
	Stop()
	GetMsgHandler() IMsgHandler
}
