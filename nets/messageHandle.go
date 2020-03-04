package nets

import (
	"fmt"
	"strconv"
	"topaza/interfaces"
)

// 消息处理模块
type MessageHandle struct {
	// 存放每一个 message id 对应的处理方法
	APIs map[uint32]interfaces.IRouter
}

// 初始化 MessageHandle
func NewMsgHandle() *MessageHandle {
	return &MessageHandle{
		APIs: make(map[uint32]interfaces.IRouter),
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
