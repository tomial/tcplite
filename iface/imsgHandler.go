package iface

// Handler抽象层

type IMsgHandler interface {
	Handle(request IRequest)

	AddRouter(msgId uint32, router IRouter)

	StartWorkerPool()

	StartWorker(int, chan IRequest)

	SendMsgToTaskQueue(IRequest)
}
