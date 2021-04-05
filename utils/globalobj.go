package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/tomial/tcplite/iface"
)

// globalobj 是全局配置对象
type GlobalObj struct {
	/*
		Server 配置
	*/

	// 全局Server对象
	TcpServer iface.IServer
	// 当前服务器监听的IP
	Host string
	// 当前服务器监听的端口号
	TcpPort int
	// 当前服务器的名称
	Name string

	/*
		其他配置信息
	*/

	// tcplite版本
	version string
	// 允许最大连接数
	MaxConn int
	// 允许最大数据包的字节数
	MaxPacketSize uint32

	// Worker池大小
	WorkerPoolSize uint32
	// 最大worker池大小
	maxWorkerPoolSize uint32
	MaxWorkerTaskLen  uint32
}

var GlobalObject *GlobalObj

// 加载用户自定义的配置文件
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/conf.json") // 读取json配置文件
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func (g *GlobalObj) GetVersion() string {
	return g.version
}

func (g *GlobalObj) GetMaxWorkerPoolSize() uint32 {
	return g.maxWorkerPoolSize
}

func init() {
	// 没有配置文件的默认值
	GlobalObject = &GlobalObj{
		Name:              "TcpLiteServerAPP",
		version:           "0.9",
		TcpPort:           8999,
		Host:              "127.0.0.1",
		MaxConn:           1000,
		MaxPacketSize:     4096,
		WorkerPoolSize:    10,
		MaxWorkerTaskLen:  1024,
		maxWorkerPoolSize: 1024,
	}

	// 尝试加载用户自定义的配置文件
	GlobalObject.Reload()

}
