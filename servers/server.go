package servers

import "hyperbird/servers/booklibrary"

// 此包用于管理各个API服务器的启动和关闭

// API基础地址
const (
	APIBasePath = "/api"
)

func Runservers() {
	// 服务器启动
	// 服务器启动

	booklibrary.RunServer()
}

func PreTestServer() {
	booklibrary.PreTestBeforeServerStart()
}

func TestServers() {
	booklibrary.Test()
}
