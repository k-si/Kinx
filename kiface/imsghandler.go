package kiface

type IMsgHandler interface {
	DoHandle(IRequest)
	InitWorkerPool()
	AllotTask(IRequest)
	GetApis() map[uint32]IRouter
}
