package main

import (
	"fmt"
	"io"
	"kinx/knet"
	"net"
	"time"
)

func main() {
	conn, _ := net.Dial("tcp4", "127.0.0.1:9999")

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
		io.ReadFull(conn, headBuf)
		msg, _ = datapack.UnPack(headBuf)

		// 根据长度读取data
		length := msg.GetMsgLen()
		dataBuf := make([]byte, length)
		io.ReadFull(conn, dataBuf)
		msg.SetMsgData(dataBuf)

		fmt.Println("id =", msg.GetMsgId(), "len =", msg.GetMsgLen(), "data =", string(msg.GetMsgData()))
	}

}
