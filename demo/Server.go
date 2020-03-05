package main

import (
	"fmt"
	"topaza/interfaces"
	"topaza/nets"
)

// 自定义路由
type PingRouter struct {
	nets.BaseRouter
}

// 在处理 conn 业务之前的方法（hook）
//func (pr *PingRouter) PreHandle(request interfaces.IRequest) {
//	fmt.Println("Call router preHandle")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
//	if err != nil {
//		fmt.Println("Call back before ping error:", err)
//	}
//}

// 在处理 conn 业务的主方法
func (pr *PingRouter) Handle(request interfaces.IRequest) {
	fmt.Println("call ping router handle")
	fmt.Println("receive from client msgId:", request.GetMsgID(),
		" data =", string(request.GetData()))

	// 读取客户端数据，再回写
	err := request.GetConnection().SendMsg(0, []byte("ping ..."))
	if err != nil {
		fmt.Println("server send message error:", err)
	}
}

// 在处理 conn 业务之后的方法（hook）
//func (pr *PingRouter) PostHandle(request interfaces.IRequest) {
//	fmt.Println("Call router postHandle")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
//	if err != nil {
//		fmt.Println("Call back after ping error:", err)
//	}
//}

type HiRouter struct {
	nets.BaseRouter
}

func (pr *HiRouter) Handle(request interfaces.IRequest) {
	fmt.Println("call hi router handle")
	fmt.Println("receive from client msgId:", request.GetMsgID(),
		" data =", string(request.GetData()))

	// 读取客户端数据，再回写
	err := request.GetConnection().SendMsg(1, []byte("hi ..."))
	if err != nil {
		fmt.Println("server send message error:", err)
	}
}

// 创建连接之后执行 hook 函数
func DoConnBegin(conn interfaces.IConnection) {
	fmt.Println("conn begin call...")
	if err := conn.SendMsg(200, []byte("conn begin")); err != nil {
		fmt.Println(err)
	}

	// 设置连接属性
	fmt.Println("set conn property:")
	conn.SetProperty("name", "ellery")

}

// 断开连接之前需要执行的函数
func DoConnEnd(conn interfaces.IConnection) {
	fmt.Println("conn end call...")
	fmt.Println("connID:", conn.GetConnID(), " lost")

	// 获取连接属性
	if name, err := conn.GetProperty("name"); err == nil {
		fmt.Println("conn property name:", name)
	}
}

// 基于框架开发的服务器端应用程序
func main() {
	// 创建一个 Server 模块
	s := nets.NewServer()

	// 注册连接 hook 函数
	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnEnd)

	// 框架添加自定义 router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HiRouter{})

	// 启动 Server
	s.Serve()
}
