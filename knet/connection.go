package knet

import (
	"fmt"
	"kinx/kiface"
	"kinx/utils"
	"net"
)

// 将server阻塞获取的连接进行封装
type Connection struct {
	conn      *net.TCPConn
	connID    uint32
	isClosed  bool
	existChan chan bool         // 客户端通知连接关闭
	router    kiface.IRouter
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.conn
}

func (c *Connection) StartReader() {
	fmt.Println("start reader goroutine:", c.connID, "remote addr:", c.conn.RemoteAddr())

	defer fmt.Println("stop reader:", c.connID, "remote addr:", c.conn.RemoteAddr())
	defer c.Stop()

	for {
		buf := make([]byte, utils.Config.MaxPackage)
		c.conn.Read(buf)

		// 处理读业务
		req := NewRequest(c, buf)
		c.router.PreHandle(req)
		c.router.Handle(req)
		c.router.PostHandle(req)
	}
}

func (c *Connection) Start() {

	// 负责从客户端读数据的业务
	go c.StartReader()

	// 负责从客户端写数据的业务
	// go c.StartWriter()
}

// 停止与客户端的连接
func (c *Connection) Stop() {
	fmt.Println("stop conn:", c.connID, "remote addr:", c.conn.RemoteAddr())

	// 去重
	if c.isClosed == true {
		return
	}

	// 停止、回收资源
	c.conn.Close()
	c.isClosed = true
	close(c.existChan)
}

func NewConnection(conn *net.TCPConn, id uint32, router kiface.IRouter) kiface.IConnection {
	c := &Connection{
		conn:      conn,
		connID:    id,
		isClosed:  false,
		existChan: make(chan bool, 1),
		router: router,
	}
	return c
}
