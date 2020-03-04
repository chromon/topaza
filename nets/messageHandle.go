package nets

import (
	"fmt"
	"strconv"
	"topaza/interfaces"
	"topaza/utils"
)

// 消息处理模块
type MessageHandle struct {
	// 存放每一个 message id 对应的处理方法
	APIs map[uint32]interfaces.IRouter

	// 负责 worker 读取任务的消息队列
	TaskQueue []chan interfaces.IRequest

	// 业务工作 worker 池
	WorkerPoolSize uint32
}

// 初始化 MessageHandle
func NewMsgHandle() *MessageHandle {
	return &MessageHandle{
		APIs: make(map[uint32]interfaces.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue: make([]chan interfaces.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// 调度执行对应的 Router 消息处理方法
func (mh *MessageHandle) DoMsgHandler(request interfaces.IRequest) {
	// 从 request 得到 msgID
	handler, ok := mh.APIs[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID:", request.GetMsgID(), " not found")
	}

	// 根据 msgID 调度对应的 router
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MessageHandle) AddRouter(msgID uint32, router interfaces.IRouter) {
	// 判断当前 message 绑定的 api 方法是否存在
	if _, ok := mh.APIs[msgID]; ok {
		// id 已经注册
		panic("repeat api, msgID:" + strconv.Itoa(int(msgID)))
	}
	// 添加
	mh.APIs[msgID] = router
	fmt.Println("add api msgID:", msgID, " success")
}

// 启动 worker 工作池，开启工作池只运行一次
func (mh *MessageHandle) StartWorkerPool() {
	// 根据 workerPoolSize 分别开启 worker
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 启动 worker
		// 开辟 worker 对应的 channel 消息队列
		mh.TaskQueue[i] = make(chan interfaces.IRequest, utils.GlobalObject.MaxWorkerTaskSize)
		// 启动 worker ，阻塞等待消息从 channel 中传递
		go mh.StartWorker(i, mh.TaskQueue[i])
	}
}

// 启动 worker 工作流程
func (mh *MessageHandle) StartWorker(workerID int, taskQueue chan interfaces.IRequest) {
	fmt.Println("Worker id:", workerID, " started...")

	// 不断阻塞等待对应消息队列的消息
	for {
		select {
			// 如果有消息过来，出列一个客户端 request，并执行当前 request 所绑定业务
			case request := <- taskQueue:
				mh.DoMsgHandler(request)
		}
	}
}

// 将消息交个 taskQueue，由 worker 处理
func (mh *MessageHandle) SendMsgToTaskQueue(request interfaces.IRequest) {
	// 将消息平均分配给不同的 worker
	// 根据客户端建立的 connID 来进行分类
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add connId:", request.GetConnection().GetConnID(),
		" request msgID:", request.GetMsgID(), " to workID:", workerID)

	// 将消息发送给对应的 taskQueue
	mh.TaskQueue[workerID] <- request
}