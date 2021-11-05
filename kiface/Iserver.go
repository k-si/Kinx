package kiface


// 服务器接口
type Iserver interface {
	Start()
	Serve()
	Stop()
}
