package kiface

type IMsgHandler interface {
	AddRouter(uint32, IRouter) IMsgHandler
	DoHandle(IRequest)
	InitWorkerPool()
	AllotTask(IRequest)
	GetApis() map[uint32]IRouter
}
