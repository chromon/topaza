package interfaces

// 封装请求消息
type IMessage interface {
	// 获取消息 id
	GetMsgId() uint32

	// 获取消息长度
	GetMsgLen() uint32

	// 获取消息内容
	GetData() []byte

	// 设置消息 id
	SetMsgId(uint32)

	// 设置消息的内容
	SetData([]byte)

	// 设置消息的长度
	SetDataLen(uint32)
}