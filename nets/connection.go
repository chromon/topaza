package nets

import (
	"errors"
	"fmt"
	"io"
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
	//// 当前连接所绑定的处理业务的方法
	//HandleAPI interfaces.HandleFunc

	// 通知当前连接已经退出/停止的 channel
	ExitChan chan bool

	// 消息的管理 MsgID 和对应的处理业务 api
	MsgHandler interfaces.IMessageHandle
}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32,
	msgHandler interfaces.IMessageHandle) *Connection {
	c := &Connection{
		Conn: conn,
		ConnID: connID,
		IsClosed: false,
		ExitChan: make(chan bool, 1),
		MsgHandler: msgHandler,
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
		// 读取客户端的数据到 buf 中， 最大 xx 字节
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("Read buf error:", err)
		//	continue
		//}

		// 创建拆包对象
		dp := NewDataPack()

		// 读取客户端 message head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read message head error:", err)
			break
		}

		// 得到 message 信息
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error:", err)
			break
		}

		// 根据 dataLen 读取 data
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read message data error:", err)
				break
			}
		}
		msg.SetData(data)

		// 得到当前 conn 的 request 请求
		req := Request {
			conn: c,
			msg: msg,
		}

		// 执行注册的路由方法从路由中，找到注册绑定的 conn 对应的 router 调用
		go c.MsgHandler.DoMsgHandler(&req)
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

// 发送数据，将数据封包后发送给远程的客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed == true {
		return errors.New("connection closed")
	}

	// 封包 |MsgDataLen|MsgId|Data|
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack error msg id: ", msgId)
		return errors.New("pack error msg")
	}

	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("write msg id: ", msgId, " error:", err)
		return errors.New("conn write error")
	}

	return nil
}