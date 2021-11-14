package knet

import (
	"fmt"
	"kinx/kiface"
	"kinx/utils"
	"net"
)

type Server struct {
	Name       string // server name
	IPVersion  string // tcp4
	IP         string
	Port       int
	MsgHandler kiface.IMsgHandler // 管理msg对应的router业务
}

func (s *Server) GetMsgHandler() kiface.IMsgHandler {
	return s.MsgHandler
}

func (s *Server) Start() {
	fmt.Println("server start:", s.Name, s.IP, s.Port, s.IPVersion)

	// 初始化worker线程池
	s.MsgHandler.InitWorkerPool()

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
		dealconn := NewConnection(conn, cid, s.MsgHandler)
		cid++
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

func NewServer() kiface.IServer {
	server := &Server{
		Name:       utils.Config.Name,
		IPVersion:  "tcp4",
		IP:         utils.Config.Host,
		Port:       utils.Config.TcpPort,
		MsgHandler: NewMsgHandler(),
	}
	return server
}
