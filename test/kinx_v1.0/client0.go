package main

import (
	"fmt"
	"github.com/k-si/Kinx/knet"
	"io"
	"net"
	"time"
)

func heartBeat(conn net.Conn) {
	timer := time.NewTimer(20 * time.Second)
	defer timer.Stop()

	OverHeartBeat:
	for {
		time.Sleep(5 * time.Second)
		datapack := knet.NewDataPack()
		msg := knet.NewMessage(100, []byte(""))
		binMsg, err := datapack.Pack(msg)
		if err != nil {
			fmt.Println("data pack err:", err)
			continue
		}
		conn.Write(binMsg)

		select {
		// 20s后退出心跳包
		case <- timer.C:
			break OverHeartBeat
		default:
		}
	}
}

func main() {
	conn, _ := net.Dial("tcp4", "127.0.0.1:9999")

	// 客户端循环发心跳包
	go heartBeat(conn)

	for {
		time.Sleep(2 * time.Second)

		// pack
		datapack := knet.NewDataPack()
		msg := knet.NewMessage(0, []byte("ping success"))
		binMsg, err := datapack.Pack(msg)
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
		msg, _ = datapack.UnPack(headBuf)

		// 根据长度读取data
		length := msg.GetMsgLen()
		dataBuf := make([]byte, length)
		_, err = io.ReadFull(conn, dataBuf)
		if err != nil {
			fmt.Println(err)
			break
		}
		msg.SetMsgData(dataBuf)

		fmt.Println("id =", msg.GetMsgId(), "len =", msg.GetMsgLen(), "data =", string(msg.GetMsgData()))
	}
}
