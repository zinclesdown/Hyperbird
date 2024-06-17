package ginserver

import "github.com/gin-gonic/gin"

var server *gin.Engine = gin.Default()

func BeforeRun() {
	Listen("/", hello)
	Listen("/api", hello)
}

func Run() {
	server.Run(":8080")
}

// 接受一个API地址和一个函数，开始监听API地址，如果有请求则调用函数。
func Listen(apipath string, f func(c *gin.Context)) {
	server.GET(apipath, f)
}

func hello(c *gin.Context) {
	c.String(200, "Hello from Hyperbird:servers:ginserver!")
}
