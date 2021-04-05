package iface

type IConnManager interface {
	// 添加连接
	Add(conn IConnection)
	// 删除连接
	Remove(conn IConnection)
	// 根据ConnID获得连接
	Get(connID uint32) (IConnection, error)
	// 获得当前连接总数
	Count() int
	// 清除并终止所有连接
	ClearAllConn()
}
