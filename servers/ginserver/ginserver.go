package ginserver

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var server *gin.Engine = gin.Default()

type RouteInfo struct {
	Path    string
	Func    func(c *gin.Context)
	Comment string
}

var routes []RouteInfo // 储存已监听的API地址和函数

func BeforeRun() {
	Listen("/", hello, "路径为根/的测试API")
	Listen("/api", hello, "路径为/api的测试API")
}

func Run(port string) {
	server.Run(":" + port)
}

// 接受一个API地址和一个函数，开始监听API地址，如果有请求则调用函数。
// 在其他地方调用它吧。
func Listen(apipath string, f func(c *gin.Context), comment string) {
	server.GET(apipath, f)
	routes = append(routes, RouteInfo{
		Path:    apipath,
		Func:    f,
		Comment: comment,
	})
}

func hello(c *gin.Context) {
	// 显示所有已注册的地址、注释
	var allroutes string
	for _, route := range routes {
		allroutes += fmt.Sprintf("Path: %-20s\t Comments: %-20s\n", route.Path, route.Comment)
	}
	c.String(200, "Welcome to server API panel!\nAvailable registered APIs:\n\n"+allroutes)

}
