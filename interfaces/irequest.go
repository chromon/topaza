package interfaces

// 将客户端请求连接信息和请求数据包装到 Request 中
type IRequest interface {
	// 得到当前请求连接
	GetConnection() IConnection

	// 得到请求的消息数据
	GetData() []byte

	// 得到请求的消息 ID
	GetMsgID() uint32
}