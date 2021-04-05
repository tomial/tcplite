package net

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/tomial/tcplite/iface"
	"github.com/tomial/tcplite/utils"
)

type Connection struct {
	// 当前连接隶属于的Server
	TCPServer iface.IServer

	// 当前连接对象(socket)
	Conn *net.TCPConn

	// 连接ID
	ConnID uint32

	// 当前连接状态
	isClosed bool

	// 告知当前连接已关闭的channel
	ExitChan chan bool

	// 无缓冲管道，用于读写Goroutine间的消息通信
	msgChan chan []byte

	// 处理该连接的Handler集合
	msgHandler iface.IMsgHandler

	// 连接属性集合
	property map[string]interface{}

	// 保护连接属性的锁
	propertyLock sync.RWMutex
}

// 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {

	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	// 添加一个连接属性
	c.property[key] = value
}

// 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {

	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("No property found")
	}
}

// 移除连接属性
func (c *Connection) RemoveProperty(key string) {

	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

func NewConnection(tcpServer iface.IServer, conn *net.TCPConn, connID uint32, handler iface.IMsgHandler) *Connection {
	c := &Connection{
		TCPServer:  tcpServer,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		msgHandler: handler,
	}

	// 将连接添加到隶属Server的ConnManager中
	c.TCPServer.GetConnManager().Add(c)
	fmt.Printf("[Connection] Added connection [%d] to connManager. Connection amount: [%d].\n", c.GetConnID(), c.TCPServer.GetConnManager().Count())

	return c
}

// 读消息Goroutine
func (c *Connection) StartReader() {
	fmt.Printf("[Reader Goroutine] now running...handling connection for id:[%d].\n", c.ConnID)
	defer fmt.Printf("[Reader Goroutine] Reader Goroutine of connection [%d] now exits.\n", c.ConnID)
	defer c.Stop()

	for {
		packer := NewDataPacker()

		// 读取Message头部
		headerBytes := make([]byte, packer.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headerBytes); err != nil {
			fmt.Println("[Reader Goroutine] Failed to read msg header from client: ", err)
			break
		}

		// 拆包得到含MsgLen 和 ID的Message对象
		msg, err := packer.Unpack(headerBytes)
		if err != nil {
			fmt.Println("[Reader Goroutine] Failed to unpack the header bytes: ", err)
			break
		}

		// 根据MsgLen读取Message的数据部分，放到msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("[Reader Goroutine] Failed to read data: ", err)
				break
			}
		}

		msg.SetData(data)

		// 得到当前连接Request对象
		req := Request{
			Conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 开启工作池机制，交给工作池处理
			c.msgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.msgHandler.Handle(&req)
		}

		// 按顺序调用Handler
	}
}

// 写消息Goroutine，将消息发送给客户端
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine] now running...")
	defer fmt.Printf("[Writer Goroutine] Writer for %s now exits.\n", c.RemoteAddr().String())

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("[Writer Goroutine] Send data error: ", err)
				return
			}
		case <-c.ExitChan:
			// Reader退出，则Writer也退出
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Printf("Connection [%d] has started.\n", c.ConnID)
	// 开启读消息Goroutine
	go c.StartReader()
	// 开启写消息Goroutine
	go c.StartWriter()

	// 调用连接开始后Hook函数
	c.TCPServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Printf("Connection [%d] stopped.\n", c.ConnID)

	// 检测是否已经关闭
	if c.isClosed == true {
		return
	}

	c.isClosed = true

	// 调用连接结束前Hook函数
	c.TCPServer.CallOnConnStop(c)

	c.Conn.Close()

	c.ExitChan <- true

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)

	// 从连接管理器中删除该连接
	c.TCPServer.GetConnManager().Remove(c)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when sending msg")
	}

	packer := NewDataPacker()

	// 打包数据
	msgBytes, err := packer.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("Failed to pack message, id = ", msgId)
		return errors.New("Pack message error")
	}

	// 将打包好的数据给写Goroutine
	c.msgChan <- msgBytes

	return nil
}
