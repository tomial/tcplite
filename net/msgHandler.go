package net

import (
	"fmt"

	"github.com/tomial/tcplite/iface"
	"github.com/tomial/tcplite/utils"
)

type MsgHandler struct {
	// 存放msgID 对应的handler方法
	handlers map[uint32]iface.IRouter

	// 负责存放消息的消息队列
	TaskQueue []chan iface.IRequest

	// 负责从消息队列中取消息进行处理的worker Goroutine数量
	WorkerPoolSize uint32
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		handlers:       make(map[uint32]iface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, // 从全局配置中获取
		TaskQueue:      make([]chan iface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (m *MsgHandler) Handle(request iface.IRequest) {
	handler, ok := m.handlers[request.GetMsgID()]
	if !ok {
		panic(fmt.Sprintf("Handler for msgId:%d not found!\n", request.GetMsgID()))
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MsgHandler) AddRouter(msgId uint32, router iface.IRouter) {
	// 判断是否重复注册
	if _, ok := m.handlers[msgId]; ok {
		// ID已经注册过了
		panic(fmt.Sprintf("msgId %d has already registered before!\n", msgId))
	}
	// 注册handler
	m.handlers[msgId] = router
	fmt.Println("Added handler for msgId: ", msgId)
}

// 启动Worker池
func (m *MsgHandler) StartWorkerPool() {

	for i := 0; i < int(m.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan iface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)

		go m.StartWorker(i, m.TaskQueue[i])
	}

}

// 启动一个Workder进行工作
func (m *MsgHandler) StartWorker(workerID int, taskQueue chan iface.IRequest) {
	fmt.Printf("Worker [%d] started.\n", workerID)

	for {
		select {
		// 处理客户端的一个请求
		case req := <-taskQueue:
			m.Handle(req)
		}
	}
}

// 将消息交给TaskQueue，由Worker进行处理
func (m *MsgHandler) SendMsgToTaskQueue(request iface.IRequest) {
	workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Printf("Send ConnID [%d] to worker [%d]\n", request.GetConnection().GetConnID(),
		workerID)

	m.TaskQueue[workerID] <- request
}
