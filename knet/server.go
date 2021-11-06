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

// 连接所绑定的业务
func callback(conn *net.TCPConn, data []byte, cnt int) error {
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("callback err:", err)
	}
	return nil
}

func (s *Server) Start() {
	fmt.Println("server start:", s.Name, s.IP, s.Port, s.IPVersion)

	// 获取server地址
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("ResolveTCPAddr err:", err)
		return
	}

	// 监听server地址
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("ListenTCP err:", err)
		return
	}

	var cid uint32
	for {
		// 阻塞的等待客户连接
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}

		// 处理客户端连接的业务
		dealconn := NewConnection(conn, cid, callback)
		cid ++
		go dealconn.Start()
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
	server := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8888,
	}
	return server
}
