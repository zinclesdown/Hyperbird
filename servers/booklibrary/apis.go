package booklibrary

// 这里是API接口的定义
// 此处的代码由且仅由前端调用，不供后端调用。

import (
	"io"

	"github.com/fatih/color"
)

// 用户API接口,供前端用户调用
type userAPI interface {
	apiLogin(username string, password string) (UserContext, error)             // 登录接口，返回用户上下文
	apiGetAllBookIds(ctx UserContext, page int, pageSize int) ([]string, error) // 获取所有书籍 ID
	apiGetBookInfo(ctx UserContext, bookid string) (Book, error)                // 获取书籍信息
	apiGetBookFile(ctx UserContext, bookid string) (FileType, io.Reader, error) // 获取书籍文件

	apiAlive() // 保持连接
}

// 管理员API接口,供前端管理员调用
type adminAPI interface {
	apiAuthenticate(jwtToken string) (UserContext, error) // 鉴权接口，返回用户上下文
	apiAddBook(ctx UserContext, book Book) error          // 添加书籍
	apiUpdateBookInfo(ctx UserContext, book Book) error   // 更新书籍信息
	apiDeleteBook(ctx UserContext, bookid string) error   // 删除书籍
}

// 定义API地址
type APIPath string

const (
	API_USER_LOGIN         APIPath = "/api/booklibrary/login"
	API_USER_GET_USER_INFO APIPath = "/api/booklibrary/getuserinfo"
	API_USER_GET_BOOK_INFO APIPath = "/api/booklibrary/getbookinfo"
	API_ALIVE              APIPath = "/api"
)

const (
	API_ADMIN_GET_ALL_BOOKS APIPath = "/api/booklibrary/admin/getallbooks"
)

// 开始监听API。在服务器完全初始化后调用
func StartListen() {
	color.Green("booklibrary::apis::StartListen()::开始监听API...")

	// // 使用GIN框架监听API
	// gin.Default().Run(":8080")

	// // servers.GinServer
	// // 监听API_ALIVE,如果有请求则返回"HTTP OK"
	// ginserver.GinServer.GET(string(API_ALIVE), func(c *gin.Context) {
	// 	c.String(200, "HTTP OK")
	// }, nil)
}

// 供前端调用的函数们，对应后端/
// 我们认为仅供前端调用的函数应该位于定义服务器的包里，不应该在其他地方调用，所以小写开头

func apiAlive() {
	// 保持连接

}
