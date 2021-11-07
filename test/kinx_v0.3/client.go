package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, _ := net.Dial("tcp4", "127.0.0.1:8888")

	for {
		time.Sleep(2 * time.Second)

		// 发
		cnt, _ := conn.Write([]byte("hello Kinx v0.3"))

		// 收
		buf := make([]byte, 1024)
		conn.Read(buf[:cnt])
		fmt.Printf("client read: %s\n", string(buf))
	}

}
