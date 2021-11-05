package knet

import (
	"fmt"
	"kinx/kiface"
	"net"
)

type Server struct {
	Name      string // server name
	IPVersion string // tcp4
	IP        string
	Port      int
}

func (s *Server) Start() {
	fmt.Printf("[%s-%s:%d/%s starting]\n", s.Name, s.IP, s.Port, s.IPVersion)

	// 获取tcp addr
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("ResolveTCPAddr err:", err)
		return
	}

	// 监听服务器地址
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("ListenTCP err:", err)
	}

	for {
		// 阻塞的等待客户连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}

		// 处理客户端连接的业务
		go func() {
			for {
				buf := make([]byte, 512)

				cnt, err := conn.Read(buf)
				if err != nil {
					fmt.Println("conn.Read err", err)
					continue
				}

				fmt.Printf("[kinx read: %s]\n", string(buf))
				fmt.Println("[kinx write]")

				_, err = conn.Write(buf[:cnt])
				if err != nil {
					fmt.Println("conn.Write err", err)
					continue
				}
			}
		}()
	}

}

func (s *Server) Serve() {
	go s.Start()

	// 处理其他扩展业务

	// 阻塞，防止主线程退出
	select {}
}

func (s *Server) Stop() {

}

func NewServer(name string) kiface.Iserver {
	server := Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8888,
	}
	return &server
}
