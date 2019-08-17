package interfaces

import "net"

// 定义连接模块的抽象层
type IConnection interface {
	// 启动连接，让当前连接准备开始工作
	Start()

	// 停止连接，结束当前连接的工作
	Stop()

	// 获取当前连接绑卡的 Socket Connection
	TCPConnection() *net.TCPConn

	// 获取当前连接模块的连接 id
	ConnID() uint32

	// 获取远程客户端的 TCP 状态 IP port
	RemoteAddr() net.Addr

	// 发送数据，讲数据发送给远程的客户端
	Send(data []byte) error
}

// 定义处理连接业务的方法
// TCP 连接，连接处理的数据内容，数据长度
type HandleFunc func(*net.TCPConn, []byte, int) error