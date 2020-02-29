package interfaces

// 定义一个服务器接口
type IServer interface {
	// 启动服务器
	Start()

	// 停止服务器
	Stop()

	// 运行服务器
	Serve()

	// 添加路由：将当前服务注册路由方法，供客户端连接处理使用
	AddRouter(router IRouter)
}