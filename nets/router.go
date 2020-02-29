package nets

import "topaza/interfaces"

// 实现 router 时，先嵌入 BaseRouter 基类，然后根据需求重写
type BaseRouter struct {}

// 在处理 conn 业务之前的方法（hook）
func (br *BaseRouter) PreHandle(request interfaces.IRequest) {}

// 在处理 conn 业务的主方法
func (br *BaseRouter) Handle(request interfaces.IRequest) {}

// 在处理 conn 业务之后的方法（hook）
func (br *BaseRouter) PostHandle(request interfaces.IRequest) {}