package interfaces

// 将请求的消息封装到 Message 中
type IMessage interface {

	// 获取消息 id
	GetMsgId() uint32

	// 设置消息 id
	SetMsgId(uint32)

	// 获取消息长度
	GetDataLen() uint32

	// 设置消息长度
	SetDataLen(uint32)

	// 获取消息的内容
	GetData() []byte

	// 设置消息的内容
	SetData([]byte)
}