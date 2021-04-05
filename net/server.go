package net

import (
	"fmt"
	"net"

	"github.com/tomial/tcplite/iface"
	"github.com/tomial/tcplite/utils"
)

// IServer 接口的实现，定义服务器模块
type Server struct {
	// 服务器名称
	Name string

	// 使用的IP版本
	IPVer string
	// 监听的IP地址
	IPAddr string

	// 监听的端口
	Port int

	// msgHandler对象
	msgHandler iface.IMsgHandler

	// 连接管理
	connManager iface.IConnManager

	// 连接开始前的Hook函数
	onConnStart func(iface.IConnection)

	// 连接开始后的Hook函数
	onConnStop func(iface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("[Server] Server Name: %s\n", utils.GlobalObject.Name)
	fmt.Printf("[Server] Starting server to listen [%s:%d]\n", utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Server] MaxConn: %d, MaxPacketSize: %d, Version: %s\n",
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize,
		utils.GlobalObject.GetVersion())
	go func() {
		s.msgHandler.StartWorkerPool()

		// 获取TCP地址
		addr, err := net.ResolveTCPAddr(s.IPVer, fmt.Sprintf("%s:%d", s.IPAddr, s.Port))
		if err != nil {
			fmt.Println("[Server] Failed to resolve TCPAddr: ", err)
		}

		// 监听目标地址
		listener, err := net.ListenTCP(s.IPVer, addr)
		if err != nil {
			fmt.Println("[Server] Failed to listen TCPAddr: ", err)
		}

		fmt.Printf("[Server] Successfully started server. Listening [%s:%d]\n", s.IPAddr, s.Port)
		var cid uint32
		cid = 0

		for {
			// 阻塞等待客户端请求，进行响应
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("[Server] Failed to accept from client: ", err)
				continue
			}

			// 判断连接是否已经超过最大连接数，超过则关闭连接
			if s.connManager.Count() >= utils.GlobalObject.MaxConn {
				// TODO 发送连接数已满的提示消息给客户端
				fmt.Println("[Server] Too many connections established.")
				conn.Close()
				continue
			}

			// 连接成功，处理请求
			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			cid++

			// 启动处理业务
			dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// 清除所有连接
	s.connManager.ClearAllConn()
	fmt.Println("[Server] server has stopped.")
}

func (s *Server) Serve() {
	// 启动服务器
	s.Start()

	// TODO 一些启动后的业务

	// 阻塞
	select {}
}

func (s *Server) AddRouter(msgId uint32, router iface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
	fmt.Println("[Server] Successfully added router.")
}

func (s *Server) GetConnManager() iface.IConnManager {
	return s.connManager
}

// 注册连接开始前Hook函数
func (s *Server) SetOnConnStart(hookFunc func(conn iface.IConnection)) {
	s.onConnStart = hookFunc
}

// 连接开始前调用Hook函数
func (s *Server) CallOnConnStart(conn iface.IConnection) {
	if s.onConnStart != nil {
		fmt.Println("[Server] Call OnConnStart().")
		s.onConnStart(conn)
	}
}

// 注册连接结束前Hook函数
func (s *Server) SetOnConnStop(hookFunc func(iface.IConnection)) {
	s.onConnStop = hookFunc
}

// 连接结束前调用Hook函数
func (s *Server) CallOnConnStop(conn iface.IConnection) {
	if s.onConnStop != nil {
		fmt.Println("[Server] Call OnConnStop().")
		s.onConnStop(conn)
	}
}

func NewServer(name string) iface.IServer {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVer:       "tcp4",
		IPAddr:      utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		msgHandler:  NewMsgHandler(),
		connManager: NewConnManager(),
	}

	return s
}
