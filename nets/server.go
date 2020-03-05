package nets

import (
	"fmt"
	"net"
	"topaza/interfaces"
	"topaza/utils"
)

// IServer 的接口实现， 定义一个 Server 的服务器模块
type Server struct {
	// 服务器名称
	Name string

	// 服务器绑定的 IP 版本
	Network string

	// 服务器监听的 IP
	IP string

	// 服务器监听的 Port
	Port int

	// 当前 server 的消息管理模块，用来绑定 msgID 和对应的处理业务 API 关系
	MsgHandler interfaces.IMessageHandle

	// server 连接管理器
	ConnManager interfaces.IConnManager

	// server 创建连接后，自动调用的 hook 函数
	OnConnStart func(conn interfaces.IConnection)

	// server 销毁连接后，自动调用的 hook 函数
	OnConnStop func(conn interfaces.IConnection)
}

// 启动服务器
func (s *Server) Start() {

	fmt.Printf("Server name: %s, version: %s\n",
		utils.GlobalObject.Name, utils.GlobalObject.Version)
	fmt.Printf("Start Server listener at IP: %s, Port: %d\n", s.IP, s.Port)

	// 服务器会一直阻塞在等待连接，应该添加到协程中
	go func() {
		// 开启消息队列及 worker 工作池
		s.MsgHandler.StartWorkerPool()

		// 获取 TCP 的 Addr
		addr, err := net.ResolveTCPAddr(s.Network, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Resolve tcp addr error:", err)
			return
		}

		// 监听服务器的地址
		listener, err := net.ListenTCP(s.Network, addr)
		if err != nil {
			fmt.Println("Listen", s.Network, " error:", err)
			return
		}
		fmt.Println("Start server success, listening...")


		// 阻塞等待客户端连接，处理客户端连接业务
		// 默认 conn ID
		var cid uint32 = 0
		for {
			// 如果有客户端连接，阻塞返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error:", err)
				continue
			}

			// 连接最大数量判断，如果超过最大连接数量，则关闭当前连接
			if s.ConnManager.Len() >= utils.GlobalObject.MaxConn {
				// TODO 个客户的响应超出最大连接信息
				fmt.Println("too many connections, max conn count:", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			// 将处理新连接的业务方法和 conn 进行绑定，得到连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			// 连接 ID 自增
			cid++

			// 启动当前连接业务处理
			go dealConn.Start()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	// 将资源、状态、连接回收停止
	fmt.Println("server stop")
	s.ConnManager.ClearConn()
}

// 运行服务器
func (s *Server) Serve() {
	// 启动 Server 的服务功能
	s.Start()

	// TODO 做一些启动服务之后的额外业务

	// 阻塞状态
	select {}
}

// 添加路由：将当前服务注册路由方法，供客户端连接处理使用
func (s *Server) AddRouter(msgID uint32, router interfaces.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add router success.")
}

func (s *Server) GetConnManager() interfaces.IConnManager {
	return s.ConnManager
}

// 初始化 Server 模块
func NewServer() interfaces.IServer {
	s := &Server{
		Name: utils.GlobalObject.Name,
		Network: "tcp4",
		IP: utils.GlobalObject.Host,
		Port: utils.GlobalObject.TCPPort,
		MsgHandler: NewMsgHandle(),
		ConnManager: NewConnManager(),
	}
	return s
}

// 注册 OnConnStart hook 函数的方法
func (s *Server) SetOnConnStart(hookFunc func(conn interfaces.IConnection)) {
	s.OnConnStart = hookFunc
}

// 注册 OnConnStop hook 函数的方法
func (s *Server) SetOnConnStop(hookFunc func(conn interfaces.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用 OnConnStart hook 函数的方法
func (s *Server) CallOnConnStart(conn interfaces.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("call OnConnStart()")
		s.OnConnStart(conn)
	}
}

// 调用 OnConnStop hook 函数的方法
func (s *Server) CallOnConnStop(conn interfaces.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("call OnConnStop()")
		s.OnConnStop(conn)
	}
}