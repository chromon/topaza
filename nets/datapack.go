package nets

import (
	"bytes"
	"encoding/binary"
	"errors"
	"topaza/interfaces"
	"topaza/utils"
)

// 封包拆包模块
type DataPack struct {}

// 初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取数据包 head
func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen uint32 (4个字节) + id uint32（4个字节）
	return 8
}

// 封包
// |dataLen|msgID|data|
func (dp *DataPack) Pack(msg interfaces.IMessage) ([]byte, error) {
	// 创建存放 bytes 字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将 dataLen 写入 dataBuff 中
	err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen())
	if err != nil {
		return nil, err
	}

	// 将 MsgID 写入 dataBuff 中
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}

	// 将 data 数据写入 dataBuff 中
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetData())
	if err != nil {
		return nil, err
	}

	// 放回二进制
	return dataBuff.Bytes(), nil
}

// 拆包，先读取 head 信息，再根据 head 信息读取
func (dp *DataPack) Unpack(binaryData [] byte) (interfaces.IMessage, error) {
	// 创建一个从输入二进制数据的 IOReader
	dataBuff := bytes.NewReader(binaryData)

	// 解压 head 信息，得到 长度和 id
	msg := &Message{}

	// 读取 dataLen
	err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}

	// 读取 MsgId
	err = binary.Read(dataBuff, binary.LittleEndian, &msg.Id)
	if err != nil {
		return nil, err
	}

	// 判断 dataLen 是否已经超出了允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large message data pack")
	}

	return msg, nil
}