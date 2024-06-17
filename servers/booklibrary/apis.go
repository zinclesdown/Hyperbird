package booklibrary

// 这里是API接口的定义
// 此处的代码由且仅由前端调用，不供后端调用。

import (
	"fmt"
	"hyperbird/servers/ginserver"
	"net/http"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

// // 用户API接口,供前端用户调用
// type userAPI interface {
// 	apiLogin(c *gin.Context)         // 登录接口，返回用户上下文
// 	apiGetAllBookIds(c *gin.Context) // 获取所有书籍 ID
// 	apiGetBookInfo(c *gin.Context)   // 获取书籍信息
// 	apiGetBookFile(c *gin.Context)   // 获取书籍文件
// 	apiAlive(c *gin.Context)         // 保持连接
// }

// // 管理员API接口,供前端管理员调用
// type adminAPI interface {
// 	apiAuthenticate(c *gin.Context)   // 鉴权接口，返回用户上下文
// 	apiAddBook(c *gin.Context)        // 添加书籍
// 	apiUpdateBookInfo(c *gin.Context) // 更新书籍信息
// 	apiDeleteBook(c *gin.Context)     // 删除书籍
// }

const ( // 用户/游客API地址
	API_USER_LOGIN               string = "/api/booklibrary/login"
	API_USER_GET_USER_INFO       string = "/api/booklibrary/getuserinfo"
	API_USER_GET_BOOK_INFO_BY_ID string = "/api/booklibrary/getbookinfobyid"
	API_USE_GET_ALL_BOOK_IDS     string = "/api/booklibrary/getallbookids"
	API_ALIVE                    string = "/api/booklibrary/"
)

const ( // 管理员/开发者API地址
	API_ADMIN_GET_ALL_BOOKS string = "/api/booklibrary/admin/getallbooks"
)

const ( // 仅供开发者调试的API地址
	APIDEV_GET_ALL_BOOKS string = "/api/booklibrary/dev/getallbooks"
)

// 注册监听相关的API。在服务器完全初始化后调用。 在server包里调用一次。
func RegisterAPIs() {
	color.Green("booklibrary::apis::RegisterAPIs()::开始注册了API...")

	ginserver.Listen(API_ALIVE, apiAlive, "图书馆系统Alive可用性测试API")
	ginserver.Listen(API_USE_GET_ALL_BOOK_IDS, apiGetAllBookIds, "获取所有保存的书籍信息")
	ginserver.Listen(API_USER_GET_BOOK_INFO_BY_ID, apiGetBookInfoById, "根据ID信息，精准获取书籍信息. 接受参数：book_id")

	// 仅供开发者调试的API
	ginserver.Listen(APIDEV_GET_ALL_BOOKS, apiDevGetAllBooks, "获取所有书籍信息")
}

//
// 以下函数供前端调用，注册于GIN服务器单例里。
//

func apiAlive(c *gin.Context) {
	// 保持连接
	c.String(200, "Alive from booklibrary!")
}

func apiGetAllBookIds(c *gin.Context) {
	// 获取所有书籍 ID
	strArr, err := GetAllBookIds(0, 10) // String[], Error
	fmt.Println("  读取书籍ID:", strArr, err)

	if err != nil {
		// 如果出错，返回错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回书籍 ID
	c.JSON(http.StatusOK, gin.H{"book_ids": strArr})
}

// 接受包内书籍ID参数，返回书籍信息 或者 错误信息
func apiGetBookInfoById(c *gin.Context) {
	// 获取书籍信息
	bookID := c.Query("book_id")         // 获取书籍ID
	book, err := GetBookInfoById(bookID) // Book, Error
	fmt.Println("  读取书籍信息:", book, err)

	if err != nil {
		// 如果出错，返回错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回书籍信息
	c.JSON(http.StatusOK, gin.H{"book": book})
}

// DEV ONLY
func apiDevGetAllBooks(c *gin.Context) {
	// 获取所有书籍信息,在终端打印他们而不是打印在网页上
	books, err := GetAllBookIds(0, 10)
	fmt.Println("  读取书籍信息:", books, err)

}
