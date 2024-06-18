package ginserver

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var server *gin.Engine = gin.Default()

type RouteInfo struct {
	Path    string
	Func    func(c *gin.Context)
	Comment string
}

var routes []RouteInfo // 储存已监听的API地址和函数

func BeforeRun(allowAllOrigins bool) {
	config := cors.DefaultConfig()
	if allowAllOrigins {
		config.AllowAllOrigins = true
	} else {
		config.AllowOrigins = []string{"http://localhost", "http://0.0.0.0", "http://127.0.0.1"}
	}
	server.Use(cors.New(config))

	Listen("/", hello, "路径为根/的测试API")
	Listen("/api", hello, "路径为/api的测试API")
}

var ListenIp, ListenPort string // 记录监听的IP和端口

// 监听IP和端口，开始运行服务器
func Run(listenip, listenport string) {
	ListenIp = listenip
	ListenPort = listenport

	server.Run(":" + ListenPort)
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
	allroutes += "<h1>Hyperbird Ginserver API Debug Panel</h1>"
	allroutes += "<p>所有已注册的可访问API均在下面的目录里。</p>"
	allroutes += "<p>监听IP: " + ListenIp + "</p>"
	allroutes += "<p>监听端口: " + ListenPort + "</p>"
	allroutes += "<table>"
	for _, route := range routes {
		allroutes += fmt.Sprintf("<tr> <td><a href='%s' target='_blank'>%s</a></td> <td>%s</td></tr>", route.Path, route.Path, route.Comment)
		// allroutes += "<br/>"
	}
	allroutes += "</table>"
	html := allroutes
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
