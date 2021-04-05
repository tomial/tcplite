package iface

import (
	"net"
)

// 封装连接套接字和处理函数
type IConnection interface {
	// 启动连接
	Start()

	// 停止连接
	Stop()

	// 获取当前连接绑定的socket对象Conn
	GetTCPConnection() *net.TCPConn

	// 获取当前连接的ID
	GetConnID() uint32

	// 获取远程客户端的TCP状态和端口
	RemoteAddr() net.Addr

	// 发送数据给客户端
	SendMsg(uint32, []byte) error

	// 设置连接属性
	SetProperty(key string, value interface{})

	// 获取连接属性
	GetProperty(key string) (interface{}, error)

	// 移除连接属性
	RemoveProperty(key string)
}

type HandleFunc func(*net.TCPConn, []byte, int) error
