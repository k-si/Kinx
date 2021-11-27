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
	"time"
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
	AcceptExitChan    chan struct{}                       // server通知accept退出
	HeartBExitChan    chan struct{}                       // server通知心跳检测退出
	FinishExitChan    chan struct{}                       // accept通知server完成退出
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
		fmt.Println("[have not registry AfterConnSuccess]")
	}
}

func (s *Server) CallBeforeConnDestroy(connection kiface.IConnection) {
	if s.BeforeConnDestroy != nil {
		s.BeforeConnDestroy(connection)
	} else {
		fmt.Println("[have not registry BeforeConnDestroy]")
	}
}

func (s *Server) AddRouter(msgId uint32, router kiface.IRouter) kiface.IServer {
	fmt.Println("[router registry SUCCESS]")

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
	fmt.Println("[server TCP start SUCCESS]:", s.Name, s.IP, s.Port, s.IPVersion)

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

OverServer:
	for {
		// 阻塞的等待客户连接
		fmt.Println("[listener accept TCP connect...]")
		conn, err := listener.AcceptTCP()
		if err != nil {
			select {
			case <-s.AcceptExitChan:
				fmt.Println("[listener receive kill signal, will clear resources]")
				s.ConnMgr.Clear()
				s.FinishExitChan <- struct{}{}
				break OverServer
			default:
			}
			fmt.Println("listener.Accept err:", err)
			continue
		}

		// 判断当前连接数量，超过max直接踢出
		if s.ConnMgr.Len() >= utils.Config.MaxConn {
			fmt.Println("[too many connections !!!]")
			conn.Close()
			continue
		}

		// 处理客户端连接的业务
		dealconn := NewConnection(s, conn, cid, s.MsgHandler)
		cid++
		go dealconn.Start()
	}
}

// TODO: 优化心跳检测算法
// 服务端心跳检测，每5s将所有连接的fresh加1
func (s *Server) heartBeat() {
	fmt.Println("[server heart beat start SUCCESS]")

OverHeartBeat:
	for {
		time.Sleep(5 * time.Second)
		for id, _ := range s.ConnMgr.GetConns() {
			conn, err := s.ConnMgr.Get(id)
			if err != nil {
				fmt.Println("get conn from commMgr err:", err)
				continue
			}
			if conn.GetFresh() == 5 {
				fmt.Println("[connection", conn.GetConnectionID(), "fresh level is 10, will stop conn!]")
				conn.Stop()
			} else {
				conn.SetFresh(conn.GetFresh() + 1)
			}
		}
		select {
		case <-s.HeartBExitChan:
			fmt.Println("[server heart beat stop]")
			break OverHeartBeat
		default:
		}
	}
}

func (s *Server) Serve() {
	fmt.Println("[server starting...]")

	// 监听系统终止进程的命令
	signal.Notify(s.DoExitChan, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 处理业务
	go s.Start()

	// 心跳检测
	go s.heartBeat()

	// 通知tcp业务goroutine退出
	<-s.DoExitChan
	s.AcceptExitChan <- struct{}{}
	s.HeartBExitChan <- struct{}{} // 这里通知heartbeat关闭的意义不大，因为heartbeat会阻塞5s，不会及时退出
	fmt.Println("[server receive kill signal, will inform listener and heart beat stop]")
	err := s.listener.Close()
	if err != nil {
		fmt.Println("close listener err:", err)
	}

	// 等待业务goroutine通知退出完成
	<-s.FinishExitChan

	// 回收资源
	s.Recycle()

	fmt.Println("[See you next time, bye~]")
}

func (s *Server) Recycle() {
	close(s.DoExitChan)
	close(s.AcceptExitChan)
	close(s.HeartBExitChan)
	close(s.FinishExitChan)
}

// 真正触发server服务停止的操作，请谨慎使用！
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
		AcceptExitChan: make(chan struct{}, 1),
		HeartBExitChan: make(chan struct{}, 1),
		FinishExitChan: make(chan struct{}, 1),
	}

	return server
}
