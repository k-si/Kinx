# Kinx
Kinx是轻量的单点serverTCP服务框架，使用多线程处理业务，读写分离，可通过配置worker线程数量或连接数量限制并发量和cpu负载。同时连接管理模块和心跳检测能很好的控制冗余资源的占用。


# 使用

### 服户端启动：

![Image text](https://ksir-oss.oss-cn-beijing.aliyuncs.com/github/kinx/kinx%E4%BD%BF%E7%94%A8.png)

### 代码中使用：

test包下有详细的使用示例，这里只做简单描述：
```go
package main

import (
	"fmt"
	"github.com/k-si/Kinx/kiface"
	"github.com/k-si/Kinx/knet"
)

type PingRouter struct {
	knet.BaseRouter
}

// 处理读业务，需要继承BaseRouter并实现Handle方法，在该方法编写具体业务逻辑。
func (p *PingRouter) Handle(req kiface.IRequest) {
	fmt.Println(">>>>>>> handler ping router")

	if err := req.GetConnection().SendMessage(0, []byte("ping client...")); err != nil {
		fmt.Println("server send message to client err:", err)
	}
}

// 其他钩子函数
func (p *PingRouter) PreHandle(req kiface.IRequest) {
	fmt.Println("pre")
}

func (p *PingRouter) PostHandle(req kiface.IRequest) {
	fmt.Println("post")
}

func main() {
	c := knet.DefaultConfig()
	s := knet.NewServer(c)

	// 业务函数注册
	pr := &PingRouter{}
	s.AddRouter(0, pr)
	
	// 启动服务
	s.Serve()
	defer s.Stop()
}
```

# 当前架构
图中的requstHandler同代码中的msgHandler
![Image text](https://ksir-oss.oss-cn-beijing.aliyuncs.com/github/kinx/Kinx0.10.png)

# 版本变更

[点此查看](https://github.com/k-si/Kinx/blob/master/version_change.md)



