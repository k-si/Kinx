package knet

import "C"
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
	exitChan   chan bool   // reader通知writer停止
	msgChan    chan []byte // reader发送writer写数据
	msgHandler kiface.IMsgHandler
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.conn
}

func (c *Connection) StartReader() {
	fmt.Println("[start reader]")

	defer fmt.Println("[stop reader]")
	defer c.Stop()

	for {
		datapack := NewDataPack()

		// 从conn中读取8个字节
		headBuf := make([]byte, MessageHeadLength)
		if _, err := io.ReadFull(c.conn, headBuf); err != nil {
			fmt.Println("read MessageHead err:", err)
			break
		}

		// 将8个字节解包成message
		msg, err := datapack.UnPack(headBuf)
		if err != nil {
			fmt.Println("unpack err:", err)
			break
		}

		// 再继续读取n个字节的data
		dataBuf := make([]byte, msg.GetMsgLen())
		if _, err := io.ReadFull(c.conn, dataBuf); err != nil {
			fmt.Println("read MessageData err:", err)
			break
		}
		msg.SetMsgData(dataBuf)

		// 将请求封装，由外部router处理读业务
		req := NewRequest(c, msg)
		// 通过消息管理器，将消息分发到对应的业务router上
		c.msgHandler.DoHandle(req)
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[start writer]")

	for {
		select {
		// 读取msgChan
		case data := <-c.msgChan:
			_, err := c.conn.Write(data)
			if err != nil {
				fmt.Println("writer send message err:", err)
			}
		// exitChan，reader通知writer关闭
		case <-c.exitChan:
			fmt.Println("[writer stop]")
			return
		}
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

	// 将二进制数据写入msgChan
	c.msgChan <- binMessage

	return nil
}

func (c *Connection) Start() {
	fmt.Println("[new connection start]", c.connID, "remote addr:", c.conn.RemoteAddr())

	// 负责从客户端读数据的业务
	go c.StartReader()

	// 负责从客户端写数据的业务
	go c.StartWriter()
}

// reader调用stop，通知writer chan
func (c *Connection) Stop() {
	fmt.Println("[stop connection]", c.connID, "remote addr:", c.conn.RemoteAddr())

	// 去重
	if c.isClosed == true {
		return
	}

	// 通知writer停止
	c.exitChan <- true

	// 回收资源
	c.conn.Close()
	close(c.exitChan)
	close(c.msgChan)
	c.isClosed = true
}

func NewConnection(conn *net.TCPConn, id uint32, msgHandler kiface.IMsgHandler) kiface.IConnection {
	c := &Connection{
		conn:       conn,
		connID:     id,
		isClosed:   false,
		exitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		msgHandler: msgHandler,
	}
	return c
}
