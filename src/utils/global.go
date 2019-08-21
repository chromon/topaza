package utils

import (
	"encoding/json"
	"io/ioutil"
	"topaza/interfaces"
)

// 存储框架全局参数

type Global struct {
	// 全局 Server 对象
	TCPServer interfaces.IServer
	// 服务器主机监听的 IP
	Host string
	// 当前服务器主机监听端口
	TCPPort int
	// 当前服务器名称
	Name string

	// 框架版本号
	Version string
	// 当前服务器允许的最大连接数
	MaxConn int
	// 当前框架运行数据包的最大值
	MaxPackageSize uint32
}

// 定义一个全局的 Global 对象
var GlobalObject *Global

// 从 json 中加载参数
func (g *Global) Reload() {
	data, err := ioutil.ReadFile("conf/global.json")
	if err != nil {
		panic(err)
	}

	// 将 json 文件数据解析到 struct 中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 提供一个 init 方法，初始化 GlobalObject 对象
func init() {

	// 如果配置文件没有加载，则使用该默认值
	GlobalObject = &Global{
		Name: "ServerApp",
		Version: "V0.4",
		TCPPort: 8989,
		Host: "0.0.0.0",
		MaxConn: 1000,
		MaxPackageSize: 4096,
	}

	// 从配置文件中加载数据
	GlobalObject.Reload()
}