package nets

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// datapack 单元测试
func TestDataPack(t *testing.T) {

	// 模拟服务端
	// 创建 socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:8989")
	if err != nil {
		fmt.Println("server listen error: ", err)
		return
	}

	// 负责从客户端处理业务
	go func() {
		// 从客户端读取数据，拆包
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error: ", err)
			}

			go func(conn net.Conn) {
				// 处理客户端的请求
				dp := NewDataPack()
				for {
					// 第一次从 conn 中读取 head
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error: ", err)
						break
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack error: ", err)
						return
					}

					if msgHead.GetDataLen() > 0 {
						// 存在数据，进行第二次从 conn 中读取内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetDataLen())

						// 根据长度再次从 io 流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data error: ", err)
							return
						}

						fmt.Println("receive msgId: ", msg.Id, ", dataLen: ",
							msg.DataLen, ", data: ", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:8989")
	if err != nil {
		fmt.Println("client dial error: ", err)
		return
	}

	dp := NewDataPack()

	// 模拟粘包，同时发送两个消息包
	msg1 := &Message {
		Id:      1,
		DataLen: 5,
		Data:    []byte{'a', 'b', 'c', 'd', 'e'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error: ", err)
		return
	}
	msg2 := &Message {
		Id:      2,
		DataLen: 5,
		Data:    []byte{'x', 'y', 'z', 'i', 'o'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error: ", err)
		return
	}

	sendData1 = append(sendData1, sendData2...)

	// 发送
	conn.Write(sendData1)

	// 客户端阻塞
	select{}
}