package nets

import "topaza/interfaces"

// 实现 IRequest 接口
type Request struct {
	// 已经和客户端建立好的连接
	conn interfaces.IConnection

	// 客户端请求数据
	data []byte
}

// 获取当前连接
func (r *Request) Connection() interfaces.IConnection {
	return r.conn
}

// 获取请求的消息数据
func (r *Request) Data() []byte {
	return r.data
}