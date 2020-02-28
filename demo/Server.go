package main

import "topaza/nets"

// 基于框架开发的服务器端应用程序
func main() {
	// 创建一个 Server 模块
	s := nets.NewServer("[Server V0.1]")
	// 启动 Server
	s.Serve()
}
