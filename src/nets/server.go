package nets

import (
	"fmt"
	"net"
	"topaza/interfaces"
)

// iServer 的接口实现，定义一个 Server 的服务器模块
type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的 IP 版本
	Network string
	// 服务器监听的 IP
	IP string
	// 服务器监听的端口
	Port int
	// 当前的 Server 添加路由
	Router interfaces.IRouter
}

// 启动服务器
func (s *Server) Start() {

	fmt.Println("[Start] Server Listener at IP: ", s.IP,
		", Port: ", s.Port, " is starting...")

	go func() {
		// IP:Port
		address := fmt.Sprintf("%s:%d", s.IP, s.Port)
		// 获取一个 TCP 的 Addr
		addr, err := net.ResolveTCPAddr(s.Network, address)
		if err != nil {
			fmt.Println("Resolve tcp addr error: ", err)
			return
		}

		// 监听服务器的地址
		listener, err := net.ListenTCP(s.Network, addr)
		if err != nil {
			fmt.Println("Listen ", s.Network, " error: ", err)
			return
		}
		fmt.Println("Start framework server", s.Name, "success, listening...")

		// 连接 id
		var cid uint32 = 0

		// 阻塞的等待客户端连接，处理客户端连接业务（读写）
		for {
			// 如果有客户端连接，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error: ", err)
				continue
			}

			// 将处理新连接的业务方法和 conn 进行绑定，得到连接模块
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			// 启动连接的业务处理
			go dealConn.Start()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	// TODO 将一些服务器资源、状态或者开辟的连接停止或回收
}

// 运行服务器
func (s *Server) Serve() {
	// 启动 Server 的服务功能
	s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞状态
	select {

	}
}

// 添加路由方法
func (s *Server) AddRouter(router interfaces.IRouter) {
	s.Router = router
	fmt.Println("Add Router success!")
}

// 初始化 Server 模块
func NewServer(name string) interfaces.IServer {
	// 创建 Server 对象
	server := &Server {
		Name: name,
		Network: "tcp4",
		IP: "0.0.0.0",
		Port: 8989,
		Router: nil,
	}

	return server
}