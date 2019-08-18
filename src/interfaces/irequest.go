package interfaces

// IRequest 接口
// 实际上是把客户端请求的连接信息和请求的数据包装到一个 Request 中

type IRequest interface {

	// 获取当前连接
	Connection() IConnection

	// 获取请求的消息数据
	Data() []byte
}
