package nets

import (
	"fmt"
	"net"
	"topaza/interfaces"
	"topaza/utils"
)

// 连接模块
type Connection struct {
	// 当前连接的 Socket TCP 套接字
	Conn *net.TCPConn

	// 连接 ID
	ConnId uint32

	// 当前的连接状态
	isClosed bool

	// 通知当前连接退出/停止 Channel
	ExitChannel chan bool

	// 连接处理方法 Router
	Router interfaces.IRouter

}

// 初始化连接模块方法
func NewConnection(conn *net.TCPConn, connId uint32,
	router interfaces.IRouter) *Connection {
	c := &Connection {
		Conn: conn,
		ConnId: connId,
		isClosed: false,
		ExitChannel: make(chan bool, 1),
		Router: router,
	}

	return c
}

// 启动连接，让当前连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Connection start, connId: ", c.ConnId)

	// 启动从当前连接的读取业务
	go c.StartReader()
	// TODO 启动从当前连接写数据业务
}

// 连接的读取业务
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")

	defer fmt.Println("ConnId: ", c.ConnId, "reader exit, remoteAddr: ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到 buf 中
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("Receive buf err: ", err)
			continue
		}

		// 得到当前 conn 数据的 Request 请求数据
		req := Request {
			conn: c,
			data: buf,
		}

		// 执行注册的路由方法
		go func(request interfaces.IRequest) {
			// 从路由中，找到注册绑定的 Conn 对应的 router 调用
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}

// 停止连接，结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Connection stop, connId: ", c.ConnId)

	// 判断当前连接是否已关闭
	if c.isClosed == true {
		return
	}

	c.isClosed = true

	// 关闭 Socket 连接
	err := c.Conn.Close()
	if err != nil {
		fmt.Println("Close connection error: ", err)
	}
	// 关闭管道
	close(c.ExitChannel)
}

// 获取当前连接绑卡的 Socket Connection
func (c *Connection) TCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接模块的连接 id
func (c *Connection) ConnID() uint32 {
	return c.ConnId
}

// 获取远程客户端的 TCP 状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据，讲数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}