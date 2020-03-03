package utils

import (
	"encoding/json"
	"io/ioutil"
	"topaza/interfaces"
)

// 框架全局参数
type Global struct {
	// 全局 Server 对象
	TCPServer interfaces.IServer

	// 服务器监听的主机 IP
	Host string

	// 服务器监听的端口号
	TCPPort int

	// 服务器名称
	Name string

	// 框架版本号
	Version string

	// 服务器允许的最大连接数
	MaxConn int

	// 框架数据包的最大值
	MaxPackageSize uint32
}

// 全局的对外对象
var GlobalObject *Global

// 从 json 中加载自定义的参数
func (g *Global) Reload() {
	data, err := ioutil.ReadFile("conf/info.json")
	if err != nil {
		panic(err)
	}
	// 将 json 文件数据解析到 struct 中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 默认初始化当前对象
func init() {
	GlobalObject = &Global{
		Name: "ServerAPP",
		Version: "v0.1",
		Host: "127.0.0.1",
		TCPPort: 8080,
		MaxConn: 1000,
		MaxPackageSize: 4096,
	}

	// 从 conf/info.json 中加载
	//GlobalObject.Reload()
}