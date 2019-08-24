package nets

import (
	"bytes"
	"encoding/binary"
	"errors"
	"topaza/interfaces"
	"topaza/utils"
)

// 封包/拆包模块
type DataPack struct {

}

// 封包/拆包初始化
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取消息头长度
func(dp *DataPack) GetHeadLen() uint32 {
	// dataLen uint32 （4个字节） + id uint32 （4个字节）
	return 8
}

// 封包
// 数据包格式：|dataLen|MsgId|data|
func(dp *DataPack) Pack(msg interfaces.IMessage) ([]byte, error) {
	// 创建一个存放字节流的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将 dataLen 写进 dataBuff 中
	err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen())
	if err != nil {
		return nil, err
	}

	// 将 MsgId 写进 dataBuff 中
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}

	// 将 data 数据写进 dataBuff 中
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetData())
	if err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包（需要将数据包 head 中的数据长度读取出来，然后再次读取数据）
func(dp *DataPack) Unpack(binaryData []byte) (interfaces.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 读取 head 信息，得到 dataLen 和 MsgId
	msg := &Message{}

	// 读取 DataLen
	err := binary.Read(dataBuff, binary.LittleEndian, &msg.DateLen)
	if err != nil {
		return nil, err
	}

	// 读取 MsgId
	err = binary.Read(dataBuff, binary.LittleEndian, &msg.Id)
	if err != nil {
		return nil, err
	}

	// 判断 dataLen 是否已经超出最大包长度
	if (utils.GlobalObject.MaxPackageSize > 0 && msg.DateLen > utils.GlobalObject.MaxPackageSize) {
		return nil, errors.New("too large message data receive")
	}

	return msg, nil
}