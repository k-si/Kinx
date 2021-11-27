# Kinx
tcp服务框架，参考Zinx。

## v0.1
实现server模块，可以正常启动一个server与client回显对话。

## v0.2
增加connection模块，将连接句柄要处理的业务抽离，server只负责生产conn，connection模块处理句柄读写业务。

## v0.3
增加router模块，通过继承一个baserouter，来自定义conn的业务函数。router传入server，server再传入connection。

## v0.4
增加config模块，用户通过json配置框架中的参数，host、ip等等。

## v0.5
将传输的数据整合为message结构体，通过自定义TLV协议(messageType|messageDataLength|messageData)解决tcp黏包问题。

## v0.6
增加消息管理模块MsgHandler，允许用户注册多个router，每个message通过id和用户自定义的业务router一一绑定。

## v0.7

完成读写分离，一个connection启用两个goroutine分别处理读写业务，读取goroutine通过channel传输数据。

## v0.8
完善消息管理模块，增加worker线程池和消息队列，每个worker带有一个队列，所有的连接业务处理均衡分配到worker上。

## v0.9
增加连接管理模块，将所有连接放入map中，方便查看、操作当前活跃的连接。增加两个hook函数，连接建立之后、连接销毁之前。

## v0.10
增加连接属性配置，方便用户给连接配置一些property。

## v0.11
通过channel实现优雅关闭 tcp server。

## v1.0
实现服务端和客户端heart beat检测，完善一些字节和控制台输出信息。

## 当前架构：
图中的requstHandler同代码中的msgHandler
![Image text](https://ksir-oss.oss-cn-beijing.aliyuncs.com/github/kinx/Kinx0.10.png)

