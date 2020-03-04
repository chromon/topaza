package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"topaza/nets"
)

// 模拟客户端
func main() {

	fmt.Println("Client0 starting...")

	// 连接远程服务器，得到连接
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Net dial error:", err)
		return
	}

	for {
		// 将消息封包
		dp := nets.NewDataPack()
		binaryMsg, err := dp.Pack(nets.NewMsgPackage(0, []byte("client test message")))
		if err != nil {
			fmt.Println("client0 pack message error:", err)
			return
		}

		_, err = conn.Write(binaryMsg)
		if err != nil {
			fmt.Println("client0 write message error:", err)
			return
		}

		// 服务器回复信息
		// 读取 tcp 流中的 head
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("client0 read head error:", err)
			break
		}

		// head 拆包
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client0 unpack error:", err)
			break
		}

		// 由 dataLen 读取 data
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*nets.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error:", err)
				return
			}

			fmt.Println("receive server msgID:", msg.Id,
				" len:", msg.DataLen, " data:", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}
}