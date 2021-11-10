package knet

import (
	"kinx/kiface"
	"strconv"
)

type MsgHandler struct {
	apis map[uint32]kiface.IRouter // message id -> router
}

func (md *MsgHandler) AddRouter(msgId uint32, router kiface.IRouter) kiface.IMsgHandler{
	// 判断不能重复注册
	if _, ok := md.apis[msgId]; ok {
		panic("repeat register router:" + strconv.Itoa(int(msgId)))
	}
	// 添加到map
	md.apis[msgId] = router

	return md
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
	return &MsgHandler{apis: make(map[uint32]kiface.IRouter)}
}
