package main

import (
	"fmt"
	"kinx/kiface"
	"kinx/knet"
)

type PingRouter struct {
	knet.BaseRouter
}

// 用户自定义的处理业务函数
func (p *PingRouter) Handle(req kiface.IRequest) {
	fmt.Println("handle==============>>>>>>>>")
	if err := req.GetConnection().SendMessage(0, []byte("ping client...")); err != nil {
		fmt.Println("server send message to client err:", err)
	}
}

func main() {
	router := &PingRouter{}
	s := knet.NewServer()
	s.AddRouter(router)
	s.Serve()
}
