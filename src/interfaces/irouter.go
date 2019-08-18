package interfaces

// 路由抽象接口
type IRouter interface {
	// 处理 conn 业务之前的 hook 方法
	PreHandle(request IRequest)

	// 处理 conn 业务的主方法
	Handle(request IRequest)

	// 处理 conn 业务之后的 hook 方法
	PostHandle(request IRequest)
}