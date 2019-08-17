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

		// 阻塞的等待客户端连接，处理客户端连接业务（读写）
		for {
			// 如果有客户端连接，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error: ", err)
				continue
			}

			// 已经与客户端建立连接，回显最大 512 字节内容
			go func() {
				for {
					buf := make([]byte, 512)
					count, err := conn.Read(buf)
					if err != nil {
						fmt.Println("Receive buf error: ", err)
						continue
					}
					fmt.Printf("Receive client buf: %s, count: %d\n", buf, count)

					// 回显
					if _, err := conn.Write(buf[0: count]); err != nil {
						fmt.Println("Write back buf error: ", err)
						continue
					}
				}
			}()
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

// 初始化 Server 模块
func NewServer(name string) interfaces.IServer {
	// 创建 Server 对象
	server := &Server {
		Name: name,
		Network: "tcp4",
		IP: "0.0.0.0",
		Port: 8989,
	}

	return server
}