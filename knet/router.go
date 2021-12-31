package knet

import "github.com/k-si/Kinx/kiface"

// 用来被继承
type BaseRouter struct{}

func (br *BaseRouter) PreHandle(req kiface.IRequest) {}

func (br *BaseRouter) Handle(req kiface.IRequest) {}

func (br *BaseRouter) PostHandle(req kiface.IRequest) {}
