package main

import (
	"fmt"
	"kinx/kiface"
	"kinx/knet"
)

type PingRouter struct {
	knet.BaseRouter
}

func (p *PingRouter) PreHandle(req kiface.IRequest) {
	fmt.Println("preHandle")
	req.GetConnection().GetTCPConnection().Write([]byte("before ping\n"))
}

func (p *PingRouter) Handle(req kiface.IRequest) {
	fmt.Println("Handle")
	req.GetConnection().GetTCPConnection().Write([]byte("ping\n"))
}

func (p *PingRouter) PostHandle(req kiface.IRequest) {
	fmt.Println("PostHandle")
	req.GetConnection().GetTCPConnection().Write([]byte("after ping\n"))
}

func main() {
	router := &PingRouter{}
	s := knet.NewServer("Kinx v0.3")
	s.AddRouter(router)
	s.Serve()
}
