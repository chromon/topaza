package nets

type Message struct {
	// 消息 id
	Id uint32

	// 消息长度
	DataLen uint32

	// 消息内容
	Data []byte
}

// 获取消息 id
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

// 设置消息 id
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

// 获取消息长度
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

// 设置消息长度
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}

// 获取消息的内容
func (m *Message) GetData() []byte {
	return m.Data
}

// 设置消息的内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}