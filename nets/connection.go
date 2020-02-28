package nets

import (
	"fmt"
	"net"
	"topaza/interfaces"
)

// 连接模块
type Connection struct {
	// 当前连接的 socket TCP 套接字
	Conn *net.TCPConn

	// 连接 ID
	ConnID uint32

	// 当前的连接状态
	IsClosed bool

	// 当前连接所绑定的处理业务的方法
	HandleAPI interfaces.HandleFunc

	// 通知当前连接已经退出/停止的 channel
	ExitChan chan bool
}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32,
	callbackAPI interfaces.HandleFunc) *Connection {
	c := &Connection{
		Conn: conn,
		ConnID: connID,
		IsClosed: false,
		HandleAPI: callbackAPI,
		ExitChan: make(chan bool, 1),
	}
	return c
}

// 连接的读业务
func (c *Connection) StartReader() {
	fmt.Println("Reader goroutine is running...")
	defer fmt.Println("ConnId =", c.ConnID,
		" reader is exit, remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到 buf 中， 最大 512 字节
		buf := make([]byte, 512)
		n, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("Read buf error:", err)
			continue
		}

		// 调用当前连接所绑定的 HandleAPI
		if err := c.HandleAPI(c.Conn, buf, n); err != nil {
			fmt.Println("ConnID =", c.ConnID, " handle error:", err)
			break
		}
	}
}

// 启动连接，让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn start, connId =", c.ConnID)

	// 启动从当前连接读取数据业务
	// 服务器端从当前连接中读业务
	go c.StartReader()
	// TODO 服务器启动从当前连接中写数据业务
}

// 停止连接，结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn stop, connId =", c.ConnID)

	if c.IsClosed == true {
		return
	}
	c.IsClosed = true

	// 关闭 socket 连接
	err := c.Conn.Close()
	if err != nil {
		fmt.Println("Conn close error:", err)
		return
	}
	// 回收资源，关闭管道
	close(c.ExitChan)
}

// 获取当前连接的绑定 socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接模块的连接 id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的 TCP 状态 IP 和 Port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据，将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}