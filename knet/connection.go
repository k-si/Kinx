package knet

import (
	"errors"
	"github.com/k-si/Kinx/kiface"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

const (
	closed = iota
	notClosed
)

// 将server阻塞获取的连接进行封装
type Connection struct {
	tcpServer    kiface.IServer // 该连接所属的server
	conn         *net.TCPConn
	connID       uint32
	isClosed     uint32
	exitChan     chan struct{} // reader通知writer停止
	msgChan      chan []byte   // reader发送writer写数据
	msgHandler   kiface.IMsgHandler
	property     map[string]interface{} // 提供用户自定义连接属性
	propertyLock *sync.RWMutex
	fresh        uint32 // 用于检测连接新鲜程度
}

func (c *Connection) SetFresh(i uint32) {
	atomic.StoreUint32(&c.fresh, i)
}

func (c *Connection) GetFresh() uint32 {
	return atomic.LoadUint32(&c.fresh)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.conn
}

func (c *Connection) GetConnectionID() uint32 {
	return c.connID
}

func (c *Connection) SetProperty(name string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[name] = value
}
func (c *Connection) GetProperty(name string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if v, ok := c.property[name]; ok {
		return v, nil
	} else {
		log.Println("property", name, "not found")
		return nil, errors.New("property not found")
	}
}

func (c *Connection) RemoveProperty(name string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if _, ok := c.property[name]; ok {
		delete(c.property, name)
	}
}

func (c *Connection) StartReader() {
	//log.Println("[reader start, belongs to connection", c.GetConnectionID(), "]")
	//defer c.conn.Close()
	defer c.Stop()

	for {
		datapack := NewDataPack()

		// 从conn中读取8个字节
		headBuf := make([]byte, MessageHeadLength)
		if _, err := io.ReadFull(c.conn, headBuf); err != nil {
			//if strings.Contains(err.Error(), "use of closed network connection") {
			//	//fmt.Println("[there happened an err: use of closed network connection, in most cases, it doesn't matter]")
			//} else {
			//	log.Println(err)
			//}
			break
		}

		// 将8个字节解包成message
		msg, err := datapack.UnPack(headBuf)
		if err != nil {
			log.Println(err)
		}

		if msg.GetMsgLen() > 0 {

			// 再继续读取n个字节的data
			dataBuf := make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.conn, dataBuf); err != nil {
				//if strings.Contains(err.Error(), "use of closed network connection") {
				//	//fmt.Println("[there happened an err: use of closed network connection, in most cases, it doesn't matter]")
				//} else {
				//	log.Println(err)
				//}
				break
			}
			msg.SetMsgData(dataBuf)

			// 将请求封装，由外部router处理读业务
			req := NewRequest(c, msg)

			// 消息管理器将task均衡分配到worker上
			c.msgHandler.AllotTask(req)
		} else {

			// 如果是只有header，表明是心跳包
			c.SetFresh(0)
		}
	}
}

func (c *Connection) StartWriter() {
	//log.Println("[writer start, belongs to connection", c.GetConnectionID(), "]")

	for {
		select {
		// 读取msgChan
		case data := <-c.msgChan:
			_, err := c.conn.Write(data)
			if err != nil {
				log.Println("writer send message err:", err)
			}
		// exitChan，reader通知writer关闭
		case <-c.exitChan:
			//log.Println("[writer stopped]")
			return
		}
	}
}

// 向客户端发送数据
func (c *Connection) SendMessage(id uint32, data []byte) error {

	// 处理conn关闭情况
	if atomic.LoadUint32(&c.isClosed) == closed {
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
	//log.Println("[new connection start]", c.connID, "remote addr:", c.conn.RemoteAddr())

	// 负责从客户端读数据的业务
	go c.StartReader()

	// 负责从客户端写数据的业务
	go c.StartWriter()

	// 连接完成后的hook
	c.tcpServer.CallAfterConnSuccess(c)
}

// 关闭连接，回收资源，删除管理池中记录的这个conn
func (c *Connection) Stop() {
	c.StopWithNotConnMgr()

	// 将连接管理池中的连接删除
	if err := c.tcpServer.GetConnMgr().Remove(c); err != nil {
		log.Println("connection remove from connMgr fail, err:", err)
	}
	//log.Println("[remove connection from connMgr, now active conn numbers:", c.tcpServer.GetConnMgr().Len(), "]")
}

// 只关闭连接，回收资源，不管连接管理池
func (c *Connection) StopWithNotConnMgr() {
	if atomic.LoadUint32(&c.isClosed) == closed {
		return
	}

	//log.Println("[stopping connection", c.connID, "remote addr:", c.conn.RemoteAddr(), "]")

	// 连接关闭之前的hook
	c.tcpServer.CallBeforeConnDestroy(c)

	// 通知writer停止
	c.exitChan <- struct{}{}
	//log.Println("[reader is closed, writer will close]")

	// 必须放在conn.close()的前面，conn的关闭除了自身退出read业务，
	// 还有可能是外界关闭，外界关闭conn后，reader业务没有停止，再次
	// 读取数据会报"使用已关闭的连接"的错误，所以使用isClosed标志连接
	// 是否已经关闭。先将标志更改，就可以避免重复close。
	atomic.StoreUint32(&c.isClosed, closed)

	// 关闭连接
	c.conn.Close()

	// 回收资源
	close(c.exitChan)
	close(c.msgChan)

	//log.Println("[connection", c.GetConnectionID(), "exit]")
}

func NewConnection(server kiface.IServer, conn *net.TCPConn, id uint32, msgHandler kiface.IMsgHandler) kiface.IConnection {
	c := &Connection{
		tcpServer:    server,
		conn:         conn,
		connID:       id,
		isClosed:     notClosed,
		exitChan:     make(chan struct{}, 1),
		msgChan:      make(chan []byte),
		msgHandler:   msgHandler,
		property:     make(map[string]interface{}),
		propertyLock: &sync.RWMutex{},
		fresh:        0,
	}

	c.tcpServer.GetConnMgr().Add(c)
	//log.Println("add connection in connMgr, active conn =", c.tcpServer.GetConnMgr().Len())

	return c
}
