package interfaces

// 封装消息，用于处理 TCP 粘包问题
type IDataPack interface {

	// 获取消息头长度
	GetHeadLen() uint32

	// 封包
	Pack(msg IMessage) ([]byte, error)

	// 拆包
	Unpack([]byte) (IMessage, error)
}