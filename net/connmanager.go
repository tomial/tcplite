package net

import (
	"errors"
	"fmt"
	"sync"

	"github.com/tomial/tcplite/iface"
)

type ConnManager struct {
	connections map[uint32]iface.IConnection // 管理的连接集合
	connLock    sync.RWMutex                 //保护连接的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]iface.IConnection),
	}
}

// 添加连接
func (c *ConnManager) Add(conn iface.IConnection) {
	// 保护共享资源，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 将conn加入到map中
	c.connections[conn.GetConnID()] = conn
	fmt.Printf("[ConnManager] Added connection [%d]\n", conn.GetConnID())
}

// 删除连接
func (c *ConnManager) Remove(conn iface.IConnection) {
	// 保护共享资源，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 删除conn
	delete(c.connections, conn.GetConnID())
	fmt.Printf("[ConnManager] Removed connection [%d]. Connection amount: [%d].\n", conn.GetConnID(), len(c.connections))
}

// 根据ConnID获得连接
func (c *ConnManager) Get(connID uint32) (iface.IConnection, error) {

	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New(
			fmt.Sprintf("[ConnManager] Connection [%d] not found.\n", connID))
	}

}

// 获得当前连接总数
func (c *ConnManager) Count() int {
	return len(c.connections)
}

// 清除并终止所有连接
func (c *ConnManager) ClearAllConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	for connID, conn := range c.connections {
		conn.Stop()
		delete(c.connections, connID)
	}

	fmt.Println("[ConnManager] All connections have been removed.")
}
