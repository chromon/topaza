package main

import (
	"fmt"
	"net"
	"time"
)

// 模拟客户端
func main() {

	fmt.Println("Client starting...")

	// 连接远程服务器，得到连接
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Net dial error:", err)
		return
	}

	for {
		// 连接调用 Write 方法，向服务器发送数据
		_, err := conn.Write([]byte("Hello server!"))
		if err != nil {
			fmt.Println("Conn write error:", err)
			return
		}

		// 接收服务器回写的数据
		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Conn read error:", err)
			return
		}

		fmt.Printf("Server call back: %s, len: %d\n", buf[:n], n)
		time.Sleep(1 * time.Second)
	}
}