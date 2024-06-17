package booklibrary

// 这里是API接口的定义
// 此处的代码由且仅由前端调用，不供后端调用。

import (
	"hyperbird/servers/ginserver"
	"io"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
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

const ( // 用户/游客API地址
	API_USER_LOGIN         string = "/api/booklibrary/login"
	API_USER_GET_USER_INFO string = "/api/booklibrary/getuserinfo"
	API_USER_GET_BOOK_INFO string = "/api/booklibrary/getbookinfo"
	API_ALIVE              string = "/api/booklibrary/"
)

const ( // 管理员/开发者API地址
	API_ADMIN_GET_ALL_BOOKS string = "/api/booklibrary/admin/getallbooks"
)

// 注册监听相关的API。在服务器完全初始化后调用。 在server包里调用一次。
func RegisterAPIs() {
	color.Green("booklibrary::apis::RegisterAPIs()::开始注册了API...")

	ginserver.Listen(API_ALIVE, apiAlive, "图书馆系统Alive可用性测试API")
}

//
// 以下函数供前端调用，注册于GIN服务器单例里。
//

func apiAlive(c *gin.Context) {
	// 保持连接
	c.String(200, "Alive from booklibrary!")
}
