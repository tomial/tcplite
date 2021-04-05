package iface

type IServer interface {
	// 启动服务器
	Start()

	// 关闭服务器
	Stop()

	// 运行服务器
	Serve()

	// 添加路由
	AddRouter(msgId uint32, router IRouter)

	// 获取连接管理器
	GetConnManager() IConnManager

	// 设置OnConnStart Hook 函数的方法
	SetOnConnStart(func(IConnection))
	// 调用OnConnStart Hook 函数的方法
	CallOnConnStart(IConnection)
	// 设置OnConnStop Hook 函数的方法
	SetOnConnStop(func(IConnection))
	// 调用OnConnStop Hook 函数的方法
	CallOnConnStop(IConnection)
}
