package nets

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 测试 dataPack 封包，拆包的单元测试
func TestDataPack(t *testing.T) {
	// 模拟服务器
	// 创建 socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	// 创建一个协程负责客户端处理业务
	go func() {
		// 从客户的读取数据，并拆包
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error:", err)
			}

			go func(conn net.Conn) {
				// 处理客户端请求
				dp := NewDataPack()
				for {
					// 拆包，读取 head
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error:", err)
						break
					}

					head, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack error:", err)
						return
					}

					if head.GetMsgLen() > 0 {
						// 存在数据
						// 根据 head dataLen 读取 data
						msg := head.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						// 根据 dataLen 的长度再次从 io 流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data error:", err)
							return
						}

						fmt.Println("receive msgId:", msg.Id, " dataLen:", msg.DataLen, " data:", string(msg.Data))
					}

				}

			}(conn)
		}
	}()


	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("client dial error:", err)
		return
	}

	// 创建一个封包对象
	dp := NewDataPack()
	// 模拟粘包，封装两个 message 发送
	msg1 := &Message{
		Id: 1,
		DataLen: 4,
		Data: []byte{'a', 'b', 'c', 'd'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg 1 error:", err)
	}

	msg2 := &Message{
		Id: 2,
		DataLen: 5,
		Data: []byte{'x', 'y', 'z', 'm', 'n'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg 2 error:", err)
	}

	// 拼接
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)

	select {}
}