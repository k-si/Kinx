package knet

import (
	"fmt"
	"kinx/kiface"
	"kinx/utils"
	"net"
	"strconv"
)

type Server struct {
	Name              string // server name
	IPVersion         string // tcp4
	IP                string
	Port              int
	MsgHandler        kiface.IMsgHandler                  // 管理msg对应的router业务
	ConnMgr           kiface.IConnectionManager           // 管理所有的连接
	AfterConnSuccess   func(connection kiface.IConnection) // 成功连接之后
	BeforeConnDestroy func(connection kiface.IConnection) // 销毁连接之前
}

func (s *Server) GetMsgHandler() kiface.IMsgHandler {
	return s.MsgHandler
}

func (s *Server) GetConnMgr() kiface.IConnectionManager {
	return s.ConnMgr
}

func (s *Server) SetAfterConnSuccess(hook func(connection kiface.IConnection)) kiface.IServer {
	s.AfterConnSuccess = hook
	return s
}

func (s *Server) SetBeforeConnDestroy(hook func(connection kiface.IConnection)) kiface.IServer {
	s.BeforeConnDestroy = hook
	return s
}

func (s *Server) CallAfterConnSuccess(connection kiface.IConnection) {
	if s.AfterConnSuccess != nil {
		s.AfterConnSuccess(connection)
	} else {
		fmt.Println("have not registry AfterConnSuccess")
	}
}

func (s *Server) CallBeforeConnDestroy(connection kiface.IConnection) {
	if s.BeforeConnDestroy != nil {
		s.BeforeConnDestroy(connection)
	} else {
		fmt.Println("have not registry BeforeConnDestroy")
	}
}

func (s *Server) AddRouter(msgId uint32, router kiface.IRouter) kiface.IServer {
	fmt.Println("router registry success")

	// 判断不能重复注册
	apis := s.MsgHandler.GetApis()
	if _, ok := apis[msgId]; ok {
		panic("repeat register router:" + strconv.Itoa(int(msgId)))
	}
	// 添加到map
	apis[msgId] = router

	return s
}

func (s *Server) Start() {
	fmt.Println("[server start]:", s.Name, s.IP, s.Port, s.IPVersion)

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

		// 判断当前连接数量，超过max直接踢出
		if s.ConnMgr.Len() >= utils.Config.MaxConn {
			fmt.Println("too many connections !!!")
			conn.Close()
			continue
		}

		// 处理客户端连接的业务
		dealconn := NewConnection(s, conn, cid, s.MsgHandler)
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
	// 清空所有连接
	s.ConnMgr.Clear()
}

func NewServer() kiface.IServer {
	server := &Server{
		Name:       utils.Config.Name,
		IPVersion:  "tcp4",
		IP:         utils.Config.Host,
		Port:       utils.Config.TcpPort,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnMgr(),
	}
	return server
}
