package interfaces

// 封包/拆包模块
type IDataPack interface {
	// 获取数据包 head
	GetHeadLen() uint32

	// 封包
	Pack(msg IMessage) ([]byte, error)

	// 拆包
	Unpack([] byte) (IMessage, error)
}