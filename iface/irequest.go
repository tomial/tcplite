package iface

// 封装请求连接和数据
type IRequest interface {
	// 获得当前请求的连接对象
	GetConnection() IConnection

	// 获得当前请求的数据
	GetData() []byte

	// 获得数据类型ID
	GetMsgID() uint32
}
