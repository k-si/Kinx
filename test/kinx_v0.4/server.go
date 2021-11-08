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
	fmt.Println("======preHandle======")
}

func (p *PingRouter) Handle(req kiface.IRequest) {
	fmt.Println("======Handle======", string(req.GetData()))
	req.GetConnection().GetTCPConnection().Write([]byte("ping"))
}

func (p *PingRouter) PostHandle(req kiface.IRequest) {
	fmt.Println("======PostHandle======")
}

func main() {
	router := &PingRouter{}
	s := knet.NewServer()
	s.AddRouter(router)
	s.Serve()
}
