package knet

import (
	"fmt"
	"kinx/kiface"
	"net"
	"time"
)

// 将server阻塞获取的连接进行封装
type Connection struct {
	Conn      *net.TCPConn
	ConnID    uint32
	isClosed  bool
	HandleApi kiface.HandleFunc // 一个conn绑定的业务方法
	ExistChan chan bool         // 客户端通知连接关闭
}

func (c *Connection) StartReader() {
	fmt.Println("start reader:", c.ConnID, "remote addr:", c.Conn.RemoteAddr())

	defer fmt.Println("stop reader:", c.ConnID, "remote addr:", c.Conn.RemoteAddr())
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		cnt, _ := c.Conn.Read(buf)
		fmt.Println("read from client:", string(buf))

		// 处理读业务
		if err := c.HandleApi(c.Conn, buf, cnt); err != nil {
			fmt.Println("handle api err:", err)
		}

		time.Sleep(3 * time.Second)
	}
}

func (c *Connection) Start() {
	fmt.Println("start conn:", c.ConnID, "remote addr:", c.Conn.RemoteAddr())

	// 负责从客户端读数据的业务
	go c.StartReader()

	// 负责从客户端写数据的业务
	// go c.StartWriter()
}

// 停止与客户端的连接
func (c *Connection) Stop() {
	fmt.Println("stop conn:", c.ConnID, "remote addr:", c.Conn.RemoteAddr())

	// 去重
	if c.isClosed == true {
		return
	}

	// 停止、回收资源
	c.Conn.Close()
	c.isClosed = true
	close(c.ExistChan)
}

func NewConnection(conn *net.TCPConn, id uint32, callback kiface.HandleFunc) kiface.Iconnection {
	c := &Connection{
		Conn:      conn,
		ConnID:    id,
		isClosed:  false,
		HandleApi: callback,
		ExistChan: make(chan bool, 1),
	}
	return c
}
