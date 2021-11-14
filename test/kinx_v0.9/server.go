package main

import (
	"fmt"
	"kinx/kiface"
	"kinx/knet"
)

type PingRouter struct {
	knet.BaseRouter
}

type HelloRouter struct {
	knet.BaseRouter
}

func (p *PingRouter) Handle(req kiface.IRequest) {
	fmt.Println(">>>>>>> handler ping router")

	if err := req.GetConnection().SendMessage(0, []byte("ping client...")); err != nil {
		fmt.Println("server send message to client err:", err)
	}
}

func (h *HelloRouter) Handle(req kiface.IRequest) {
	fmt.Println(">>>>>>> handler hello router")

	if err := req.GetConnection().SendMessage(0, []byte("hello client...")); err != nil {
		fmt.Println("server send message to client err:", err)
	}
}

func after(connection kiface.IConnection) {
	connection.SendMessage(200, []byte("上线啦～"))
}

func before(connection kiface.IConnection) {
	fmt.Println("下线啦～", connection.GetConnectionID())
}

func main() {
	pr := &PingRouter{}
	hr := &HelloRouter{}
	s := knet.NewServer()
	s.SetAfterConnSuccess(after).SetBeforeConnDestroy(before)
	s.AddRouter(0, pr).AddRouter(1, hr)
	s.Serve()
}
