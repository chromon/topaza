package nets

import "topaza/interfaces"

// 实现 router 时， 先实现 BaseRouter 基类，然后根据需求进行重写
type BaseRouter struct {}


// 处理 conn 业务之前的 hook 方法
func (br *BaseRouter) PreHandle(request interfaces.IRequest) {}

// 处理 conn 业务的主方法
func (br *BaseRouter) Handle(request interfaces.IRequest) {}

// 处理 conn 业务之后的 hook 方法
func (br *BaseRouter) PostHandle(request interfaces.IRequest) {}