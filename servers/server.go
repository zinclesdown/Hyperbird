package servers

// 此包用于管理各个API服务器的启动和关闭

import (
	"hyperbird/servers/booklibrary"
	"hyperbird/servers/ginserver"
)

// 开始监听所有服务器的API。在服务器完全初始化后调用
func StartAllListen() {
	ginserver.BeforeRun()

	// 注册API
	booklibrary.RegisterAPIs()
	// xxx.RegisterAPIs()
	// xxx.RegisterAPIs()
	// xxx.RegisterAPIs()

	ginserver.Run("8080") // 暂定使用8080端口，后面可以改吧
}

// 服务器启动前的初始化
func InitServers() {
	booklibrary.InitServer()

	StartAllListen()
}

// 清理目录
func PreTestServer() {
	booklibrary.PreTestBeforeServerStart()
}

func TestServers() {
	booklibrary.Test()
}
