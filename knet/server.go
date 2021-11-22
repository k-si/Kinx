package knet

import (
	"fmt"
	"kinx/kiface"
	"kinx/utils"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

type Server struct {
	Name              string // server name
	IPVersion         string // tcp4
	IP                string
	Port              int
	mu                sync.Mutex                          // 每次生成连接时需要加锁
	MsgHandler        kiface.IMsgHandler                  // 管理msg对应的router业务
	ConnMgr           kiface.IConnectionManager           // 管理所有的连接
	DoExitChan        chan os.Signal                      // 系统通知server退出
	AcceptExitChan    chan bool                           // server通知accept退出
	FinishExitChan    chan bool                           // accept通知server完成退出
	listener          *net.TCPListener                    // 监听句柄
	AfterConnSuccess  func(connection kiface.IConnection) // 成功连接之后
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
	s.listener = listener
	if err != nil {
		fmt.Println("ListenTCP err:", err)
		return
	}

	var cid uint32

	for {
		// 加锁防止并发连接请求到达
		//s.mu.Lock()

		// 阻塞的等待客户连接
		conn, err := listener.AcceptTCP()
		if err != nil {
			select {
			default:
			case <-s.AcceptExitChan:
				s.ConnMgr.Clear()
				s.FinishExitChan <- true
				return
			}
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

		//s.mu.Unlock()
	}
}

func (s *Server) Serve() {
	// 监听系统终止进程的命令
	signal.Notify(s.DoExitChan, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 处理业务
	go s.Start()

	// 通知业务goroutine退出
	<-s.DoExitChan
	s.AcceptExitChan <- true
	s.listener.Close()

	// 等待业务goroutine通知退出完成
	<-s.FinishExitChan

	// 回收资源
	s.Recycle()

	fmt.Println("Exit Kinx...bye!")
}

func (s *Server) Recycle() {
	close(s.DoExitChan)
	close(s.AcceptExitChan)
	close(s.FinishExitChan)
}

// 真正出发server服务停止的操作，请谨慎使用！
func (s *Server) Stop() {
	s.DoExitChan <- os.Kill
}

func NewServer() kiface.IServer {
	server := &Server{
		Name:           utils.Config.Name,
		IPVersion:      "tcp4",
		IP:             utils.Config.Host,
		Port:           utils.Config.TcpPort,
		MsgHandler:     NewMsgHandler(),
		ConnMgr:        NewConnMgr(),
		DoExitChan:     make(chan os.Signal, 1),
		AcceptExitChan: make(chan bool, 1),
		FinishExitChan: make(chan bool, 1),
	}

	return server
}
