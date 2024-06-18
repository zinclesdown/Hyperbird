package booklibrary

// 这里是API接口的定义
// 此处的代码由且仅由前端调用，不供后端调用。

import (
	"fmt"
	"hyperbird/servers/ginserver"
	"net/http"
	"strconv"

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
	// API_USER_LOGIN               string = "/api/booklibrary/login"
	// API_USER_GET_USER_INFO       string = "/api/booklibrary/getuserinfo"
	API_USER_GET_BOOK_INFO_BY_ID string = "/api/book_library/get_book_info_by_id"
	API_USE_GET_ALL_BOOK_IDS     string = "/api/book_library/get_all_bookids"
	API_ALIVE                    string = "/api/book_library/"

	API_GET_BOOKS_SHORT_INFO string = "/api/book_library/get_books_short_info"
)

const ( // 管理员/开发者API地址
	API_ADMIN_GET_ALL_BOOKS string = "/api/booklibrary/admin/getallbooks"
)

// 注册监听相关的API。在服务器完全初始化后调用。 在server包里调用一次。
func RegisterAPIs() {
	color.Green("booklibrary::apis::RegisterAPIs()::开始注册了API...")

	ginserver.Listen(API_ALIVE, apiAlive, "图书馆系统Alive可用性测试API")
	ginserver.Listen(API_USE_GET_ALL_BOOK_IDS, apiGetAllBookIds, "获取所有保存的书籍信息")
	ginserver.Listen(API_USER_GET_BOOK_INFO_BY_ID, apiGetBookInfoById, "根据ID信息，精准获取书籍信息. 接受参数：book_id")

	//apiGetBooksShortInfo
	ginserver.Listen(API_GET_BOOKS_SHORT_INFO, apiGetBooksShortInfo, "返回某一页的书籍ID与书籍的名称、图像。接受参数：page:int, page_size:int.")
	// 仅供开发者调试的API
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

// 返回某一页的书籍ID与书籍的名称、图像。接受参数：page:int, page_size:int.
// 返回： [{book_id:string, book_name:string, book_image_path:sting(url)}]
func apiGetBooksShortInfo(c *gin.Context) {
	// 从url参数读取page:int page_size:int
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数page错误"})
		return
	}
	page_size, err := strconv.Atoi(c.Query("page_size"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数page_size错误"})
		return
	}

	// 调用GetAllBookIds, 获取所有ID, 然后根据ID进一步索引得到书籍信息
	// 获取所有书籍 ID
	strArr, err := GetAllBookIds(page, page_size) // String[], Error
	fmt.Println("  读取书籍ID:", strArr, err)
	if err != nil {
		// 如果出错，返回错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 根据已有ID数组，获取书籍信息。
	// 定义结构：
	type bookShortInfo struct {
		BookId        string `json:"book_id"`
		BookName      string `json:"book_name"`
		BookImagePath string `json:"book_image_path"`
	}

	// 定义数组(切片)：
	bookShortInfos := make([]bookShortInfo, len(strArr))

	// 遍历ID数组，获取书籍信息
	for i, book_id := range strArr {
		book, err := GetBookInfoById(book_id) // Book, Error
		if err != nil {
			// 如果出错，返回错误信息
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 赋值
		bookShortInfos[i].BookId = book_id
		bookShortInfos[i].BookName = book.BookName
		bookShortInfos[i].BookImagePath = book.BookImagePath
	}

	// 返回这个数组
	c.JSON(http.StatusOK, gin.H{"books": bookShortInfos})
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
