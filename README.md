# Kinx
tcp服务框架，参考Zinx

## v0.1
实现server模块，可以正常启动一个server与client回显对话。

## v0.2
增加connection模块，将连接句柄要处理的业务抽离，server只负责生产conn，connection模块处理句柄读写业务。

## v0.3
增加router模块，通过继承一个baserouter，来自定义conn的业务函数。router传入server，server再传入connection。

## v0.4
增加config模块，用户通过json配置框架中的参数，host、ip等等


