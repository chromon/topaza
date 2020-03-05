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
	AddRouter(msgID uint32, router IRouter)

	// 获取连接管理器
	GetConnManager() IConnManager

	// 注册 OnConnStart hook 函数的方法
	SetOnConnStart(func(conn IConnection))

	// 注册 OnConnStop hook 函数的方法
	SetOnConnStop(func(conn IConnection))

	// 调用 OnConnStart hook 函数的方法
	CallOnConnStart(conn IConnection)

	// 调用 OnConnStop hook 函数的方法
	CallOnConnStop(conn IConnection)
}