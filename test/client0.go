package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

type Message struct {
	id     uint32
	length uint32
	data   []byte
}

func heartBeat(conn net.Conn) {
	timer := time.NewTimer(20 * time.Second)
	defer timer.Stop()

OverHeartBeat:
	for {
		time.Sleep(5 * time.Second)
		binMsg, err := Pack(100, []byte(""))
		if err != nil {
			fmt.Println("data pack err:", err)
			continue
		}
		conn.Write(binMsg)

		select {
		// 20s后退出心跳包
		case <-timer.C:
			break OverHeartBeat
		default:
		}
	}
}

func main() {
	conn, _ := net.Dial("tcp4", "127.0.0.1:8088")

	// 客户端循环发心跳包
	go heartBeat(conn)

	for {
		time.Sleep(2 * time.Second)

		// pack
		binMsg, err := Pack(0, []byte("ping success"))
		if err != nil {
			fmt.Println("data pack err:", err)
			continue
		}
		conn.Write(binMsg)

		// unpack
		// 读取8字节
		headBuf := make([]byte, 8)
		_, err = io.ReadFull(conn, headBuf)
		if err != nil {
			fmt.Println(err)
			break
		}
		msg, err := UnPack(headBuf)
		if err != nil {
			fmt.Println(err)
		}

		// 根据长度读取data
		length := msg.length
		dataBuf := make([]byte, length)
		_, err = io.ReadFull(conn, dataBuf)
		if err != nil {
			fmt.Println(err)
			break
		}
		msg.data = dataBuf

		fmt.Println("id =", msg.id, "len =", msg.length, "data =", string(msg.data))
	}
}

// 将message转为二进制切片
func Pack(id uint32, data []byte) ([]byte, error) {

	// 创建二进制buf，存放message
	buf := bytes.NewBuffer([]byte{})

	// 将msg中的id、len、data放入buf
	if err := binary.Write(buf, binary.LittleEndian, id); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, uint32(len(data))); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 将二进制数据中的head抽离出来
func UnPack(data []byte) (*Message, error) {

	// 将二进制数据存入buf
	buf := bytes.NewBuffer(data)
	head := &Message{}

	// 将二进制数据读取出来，存入结构体中
	if err := binary.Read(buf, binary.LittleEndian, &head.id); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &head.length); err != nil {
		return nil, err
	}

	return head, nil
}
