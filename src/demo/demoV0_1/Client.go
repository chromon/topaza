package main

import (
	"fmt"
	"net"
	"time"
)

// 模拟客户端
func main() {

	fmt.Println("Client start ...")

	time.Sleep(1 * time.Second)

	// 连接远程服务器，得到一个 conn 连接
	conn, err := net.Dial("tcp", "127.0.0.1:8989")
	if err != nil {
		fmt.Println("Client start error: ", err, ", exit!")
		return
	}

	for {
		// 连接调用 write 写数据
		_, err := conn.Write([]byte("hello framework V0.1"))
		if err != nil {
			fmt.Println("Write conn error: ", err)
		}

		buf := make([]byte, 512)
		count, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read buf error: ", err)
			return
		}

		fmt.Printf("Server call back: %s, count: %d\n", buf, count)

		// cpu 阻塞
		time.Sleep(1 * time.Second)
	}

}
