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
func (pr *PingRouter) PreHandle(request interfaces.IRequest) {
	fmt.Println("Call router preHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("Call back before ping error:", err)
	}
}

// 在处理 conn 业务的主方法
func (pr *PingRouter) Handle(request interfaces.IRequest) {
	fmt.Println("Call router handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...\n"))
	if err != nil {
		fmt.Println("Call back ping error:", err)
	}
}

// 在处理 conn 业务之后的方法（hook）
func (pr *PingRouter) PostHandle(request interfaces.IRequest) {
	fmt.Println("Call router postHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("Call back after ping error:", err)
	}
}

// 基于框架开发的服务器端应用程序
func main() {
	// 创建一个 Server 模块
	s := nets.NewServer("[Server V0.1]")
	// 框架添加自定义 router
	s.AddRouter(&PingRouter{})
	// 启动 Server
	s.Serve()
}
