package interfaces

import "net"

// 定义连接模块的抽象层
type IConnection interface {
	// 启动连接，让当前的连接准备开始工作
	Start()

	// 停止连接，结束当前连接的工作
	Stop()

	// 获取当前连接的绑定 socket conn
	GetTCPConnection() *net.TCPConn

	// 获取当前连接模块的连接 id
	GetConnID() uint32

	// 获取远程客户端的 TCP 状态 IP 和 Port
	RemoteAddr() *net.Addr

	// 发送数据，将数据发送给远程的客户端
	Send(data []byte) error
}

// 处理连接业务的方法
// *net.TCPConn 客户端连接
// []byte 处理数据内容
// int 处理数据长度
type HandleFunc func(*net.TCPConn, []byte, int) error