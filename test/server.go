package main

import (
	"fmt"
	"github.com/k-si/Kinx/kiface"
	"github.com/k-si/Kinx/knet"
)

type PingRouter struct {
	knet.BaseRouter
}

type HelloRouter struct {
	knet.BaseRouter
}

// 继承BaseRouter，编写具体业务代码
func (p *PingRouter) Handle(req kiface.IRequest) {
	fmt.Println(">>>>>>> handler ping router")
	fmt.Println(req.GetMsg())

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

// 关于connection的钩子函数
func after(connection kiface.IConnection) {
	connection.SendMessage(200, []byte(">>>>>>> 上线啦～"))

	// 给connection提供一些属性
	connection.SetProperty("name", "zhang san")
	connection.SetProperty("age", 18)
}

func before(connection kiface.IConnection) {
	fmt.Println(">>>>>>> 下线啦～", connection.GetConnectionID())
	name, _ := connection.GetProperty("name")
	age, _ := connection.GetProperty("age")
	fmt.Println(name, age)
}

func main() {
	pr := &PingRouter{}
	hr := &HelloRouter{}
	c := knet.DefaultConfig()
	c.MaxConnSize = 2
	s := knet.NewServer(c)

	s.SetAfterConnSuccess(after).SetBeforeConnDestroy(before)
	s.AddRouter(0, pr).AddRouter(1, hr).Serve()
}
