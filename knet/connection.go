package knet

import (
	"errors"
	"fmt"
	"io"
	"kinx/kiface"
	"net"
)

// 将server阻塞获取的连接进行封装
type Connection struct {
	conn       *net.TCPConn
	connID     uint32
	isClosed   bool
	existChan  chan bool // 客户端通知连接关闭
	msgHandler kiface.IMsgHandler
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.conn
}

func (c *Connection) StartReader() {
	fmt.Println("start reader goroutine:", c.connID, "remote addr:", c.conn.RemoteAddr())

	defer fmt.Println("stop reader:", c.connID, "remote addr:", c.conn.RemoteAddr())
	defer c.Stop()

	for {
		datapack := NewDataPack()

		// 从conn中读取8个字节
		headBuf := make([]byte, MessageHeadLength)
		if _, err := io.ReadFull(c.conn, headBuf); err != nil {
			fmt.Println("read MessageHead err:", err)
			continue
		}

		// 将8个字节解包成message
		msg, err := datapack.UnPack(headBuf)
		if err != nil {
			fmt.Println("unpack err:", err)
		}

		// 再继续读取n个字节的data
		dataBuf := make([]byte, msg.GetMsgLen())
		if _, err := io.ReadFull(c.conn, dataBuf); err != nil {
			fmt.Println("read MessageData err:", err)
			continue
		}
		msg.SetMsgData(dataBuf)

		// 将请求封装，由外部router处理读业务
		req := NewRequest(c, msg)
		// 通过消息管理器，将消息分发到对应的业务router上
		c.msgHandler.DoHandle(req)
	}
}

// 向客户端发送数据
func (c *Connection) SendMessage(id uint32, data []byte) error {
	// 处理conn关闭情况
	if c.isClosed {
		return errors.New("connection have been closed")
	}

	datapack := NewDataPack()

	// 将message封包成二进制数据
	binMessage, err := datapack.Pack(NewMessage(id, data))
	if err != nil {
		return errors.New("pack message failed")
	}

	// 句柄写出二进制数据
	if _, err := c.conn.Write(binMessage); err != nil {
		return errors.New("write message to client failed")
	}

	return nil
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

func NewConnection(conn *net.TCPConn, id uint32, msgHandler kiface.IMsgHandler) kiface.IConnection {
	c := &Connection{
		conn:       conn,
		connID:     id,
		isClosed:   false,
		existChan:  make(chan bool, 1),
		msgHandler: msgHandler,
	}
	return c
}
