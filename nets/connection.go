package nets

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"topaza/interfaces"
	"topaza/utils"
)

// 连接模块
type Connection struct {
	// 当前 conn 所属 server
	TCPServer interfaces.IServer

	// 当前连接的 socket TCP 套接字
	Conn *net.TCPConn

	// 连接 ID
	ConnID uint32

	// 当前的连接状态
	IsClosed bool
	//// 当前连接所绑定的处理业务的方法
	//HandleAPI interfaces.HandleFunc

	// 通知当前连接已经退出/停止的 channel，有 reader 告知 writer
	ExitChan chan bool

	// 无缓冲管道，用于读写 goroutine 之间的消息通信
	MsgChan chan []byte

	// 消息的管理 MsgID 和对应的处理业务 api
	MsgHandler interfaces.IMessageHandle

	// 连接属性集合
	property map[string]interface{}

	// 保护连接属性锁
	propertyLock sync.RWMutex
}

// 初始化连接模块的方法
func NewConnection(server interfaces.IServer, conn *net.TCPConn, connID uint32,
	msgHandler interfaces.IMessageHandle) *Connection {
	c := &Connection{
		TCPServer: server,
		Conn: conn,
		ConnID: connID,
		IsClosed: false,
		ExitChan: make(chan bool, 1),
		MsgChan: make(chan []byte),
		MsgHandler: msgHandler,
		property: make(map[string]interface{}),
	}

	// 将 conn 加入 ConnManager 中
	c.TCPServer.GetConnManager().Add(c)

	return c
}

// 连接的读业务
func (c *Connection) StartReader() {
	fmt.Println("Reader goroutine is running...")
	defer fmt.Println("reader is exit, ConnId:", c.ConnID,
		" remote addr:", c.RemoteAddr().String())
	defer c.Stop()

	for {
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

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池，将消息发送给 worker 工作池即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 执行注册的路由方法从路由中，找到注册绑定的 conn 对应的 router 调用
			go c.MsgHandler.DoMsgHandler(&req)
		}


	}
}

// 连接的写消息 goroutine，用于给客户端写消息
func (c *Connection) StartWriter() {
	fmt.Println("Writer goroutine is running...")
	defer fmt.Println("conn writer exit, ConnId:", c.ConnID, " remote addr:", c.RemoteAddr().String())

	// 阻塞等待 channel 的消息，写给客户端
	for {
		select {
		case data := <- c.MsgChan:
			// 有数据需要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error:", err)
				return
			}
		case <- c.ExitChan:
			// reader 已退出
			return
		}
	}
}

// 启动连接，让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn start, connId =", c.ConnID)

	// 启动从当前连接读取数据业务
	// 服务器端从当前连接中读业务
	go c.StartReader()
	// 服务器启动从当前连接中写数据业务
	go c.StartWriter()

	// 按照开发者传进来的创建连接后需要调用的处理业务，执行 hook 函数
	c.TCPServer.CallOnConnStart(c)
}

// 停止连接，结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn stop, connId =", c.ConnID)

	if c.IsClosed == true {
		return
	}
	c.IsClosed = true

	// 按照开发者传进来的销毁连接前需要调用的处理业务，执行 hook 函数
	c.TCPServer.CallOnConnStop(c)

	// 关闭 socket 连接
	err := c.Conn.Close()
	if err != nil {
		fmt.Println("Conn close error:", err)
		return
	}

	// 关闭 writer
	c.ExitChan <- true

	// 将当前连接从 connManager 中移除
	c.TCPServer.GetConnManager().Remove(c)

	// 回收资源，关闭管道
	close(c.ExitChan)
	close(c.MsgChan)
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

	// 将消息发送给管道
	c.MsgChan <- binaryMsg

	return nil
}

// 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	// 加写锁
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	// 添加一个连接属性
	c.property[key] = value
}

// 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	// 加读锁
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	// 读取属性
	if value, ok := c.property[key]; ok {
		return value, nil
	}

	return nil, errors.New("no property found")
}

// 删除连接属性
func (c *Connection) RemoveProperty(key string) {
	// 加写锁
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	// 删除属性
	delete(c.property, key)
}