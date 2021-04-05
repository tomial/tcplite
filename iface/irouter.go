package iface

type IRouter interface {
	// 处理业务之前的钩子方法
	PreHandle(request IRequest)
	// 处理业务的主方法
	Handle(request IRequest)
	// 处理业务之后的钩子方法
	PostHandle(request IRequest)
}
