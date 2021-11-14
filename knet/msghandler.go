package knet

import (
	"fmt"
	"kinx/kiface"
	"kinx/utils"
	"strconv"
)

type MsgHandler struct {
	apis            map[uint32]kiface.IRouter // message id -> router
	workerTaskQueue []chan kiface.IRequest    // worker任务队列
	workerPoolSize  uint32                    // worker线程池
}

func (md *MsgHandler) GetApis() map[uint32]kiface.IRouter {
	return md.apis
}

func (md *MsgHandler) AddRouter(msgId uint32, router kiface.IRouter) kiface.IMsgHandler {
	fmt.Println("router registry success")

	// 判断不能重复注册
	if _, ok := md.apis[msgId]; ok {
		panic("repeat register router:" + strconv.Itoa(int(msgId)))
	}
	// 添加到map
	md.apis[msgId] = router

	return md
}

func (md *MsgHandler) InitWorkerPool() {
	fmt.Println("init worker pool, size:", utils.Config.WorkerPoolSize)

	for i := 0; i < int(md.workerPoolSize); i++ {
		// 千万不要忘了初始化channel，虽然channel列表已经make，但其中的每一个channel也要make
		md.workerTaskQueue[i] = make(chan kiface.IRequest, utils.Config.MaxWorkerTaskSize)
		go md.StartWorker(i, md.workerTaskQueue[i])
	}
}

func (md *MsgHandler) StartWorker(workerID int, taskQueue chan kiface.IRequest) {
	fmt.Println("worker", workerID, "start!")

	// worker处理队列中的任务
	for {
		select {
		case req := <-taskQueue:
			md.DoHandle(req) // 交给router处理request
		}
	}
}

func (md *MsgHandler) AllotTask(req kiface.IRequest) {
	// connection 将request均衡的分配给每个worker
	workerId := req.GetConnection().GetConnectionID() % md.workerPoolSize

	// 按照连接id进行分配，一个连接专属一个worker
	md.workerTaskQueue[workerId] <- req

	fmt.Println("request", req.GetConnection().GetConnectionID(), "had allot to worker", workerId)
}

func (md *MsgHandler) DoHandle(req kiface.IRequest) {
	// 判断是否存在对应的router
	if _, ok := md.apis[req.GetMsg().GetMsgId()]; !ok {
		panic("no router bind to this message id:" + strconv.Itoa(int(req.GetMsg().GetMsgId())))
	}
	// 调用router处理函数
	router := md.apis[req.GetMsg().GetMsgId()]
	router.PreHandle(req)
	router.Handle(req)
	router.PostHandle(req)
}

func NewMsgHandler() kiface.IMsgHandler {
	return &MsgHandler{
		apis:            make(map[uint32]kiface.IRouter),
		workerTaskQueue: make([]chan kiface.IRequest, utils.Config.WorkerPoolSize),
		workerPoolSize:  utils.Config.WorkerPoolSize,
	}
}
