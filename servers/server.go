package servers

// 此包用于管理各个API服务器的启动和关闭

import (
	"hyperbird/servers/booklibrary"
	"hyperbird/servers/ginserver"
)

// 开始监听所有服务器的API。在服务器完全初始化后调用
func StartAllListen() {
	ginserver.BeforeRun()

	ginserver.Run()
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
