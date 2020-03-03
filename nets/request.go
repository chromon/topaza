package nets

import "topaza/interfaces"

type Request struct {
	// 已经和客户端建立好的连接
	conn interfaces.IConnection

	// 客户端请求的数据
	msg interfaces.IMessage
}

// 得到当前连接
func (r *Request) GetConnection() interfaces.IConnection {
	return r.conn
}

// 得到请求的消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}