package nets

import (
	"fmt"
	"net"
	"topaza/interfaces"
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

	// 当前 Server 添加 router
	Router interfaces.IRouter
}

// 启动服务器
func (s *Server) Start() {

	fmt.Printf("Start Server listener at IP: %s, Port: %d\n", s.IP, s.Port)

	// 服务器会一直阻塞在等待连接，应该添加到协程中
	go func() {
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
		fmt.Println("Start server success,", s.Name, "listening...")


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

			// 将处理新连接的业务方法和 conn 进行绑定，得到连接模块
			dealConn := NewConnection(conn, cid, s.Router)
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
func (s *Server) AddRouter(router interfaces.IRouter) {
	s.Router = router
	fmt.Println("Add router success.")
}

// 初始化 Server 模块
func NewServer(name string) interfaces.IServer {
	s := &Server{
		Name: name,
		Network: "tcp4",
		IP: "127.0.0.1",
		Port: 8080,
		Router: nil,
	}
	return s
}
