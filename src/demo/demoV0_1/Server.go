package main

import "topaza/nets"

// 基于框架开发的服务器应用程序

func main() {
	// 创建一个 Server，使用框架的 API
	server := nets.NewServer("[Topaza V0.1]")
	// 启动 server
	server.Serve()
}
