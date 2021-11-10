package kiface

type IMsgHandler interface {
	AddRouter(uint32, IRouter) IMsgHandler
	DoHandle(IRequest)
}
