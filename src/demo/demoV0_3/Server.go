package main

import (
	"fmt"
	"topaza/interfaces"
	"topaza/nets"
)

// 基于框架开发的服务器应用程序

// ping test 自定义路由
type PingRouter struct {
	nets.BaseRouter
}

// test PreHandle
func (pr *PingRouter) PreHandle(request interfaces.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.Connection().TCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("Call back before ping error: ", err)
	}
}

// test Handle
func (pr *PingRouter) Handle(request interfaces.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.Connection().TCPConnection().Write([]byte("ping...\n"))
	if err != nil {
		fmt.Println("Call back ping error: ", err)
	}
}

// test PostHandle
func (pr *PingRouter) PostHandle(request interfaces.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.Connection().TCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("Call back after ping error: ", err)
	}
}

func main() {
	// 创建一个 Server，使用框架的 API
	server := nets.NewServer("[Topaza V0.3]")
	// 添加自定义的路由
	server.AddRouter(&PingRouter{})
	// 启动 server
	server.Serve()
}
