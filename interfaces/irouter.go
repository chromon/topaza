package interfaces

// 路由抽象接口
type IRouter interface {
	// 在处理 conn 业务之前的方法（hook）
	PreHandle(request IRequest)

	// 在处理 conn 业务的主方法
	Handle(request IRequest)

	// 在处理 conn 业务之后的方法（hook）
	PostHandle(request IRequest)
}