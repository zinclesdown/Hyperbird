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

	API_GET_BOOKS_SHORT_INFO string = "/api/book_library/get_books_short_info" // 获取供列表界面预览的短信息
	API_GET_BOOKS_INFO       string = "/api/book_library/get_books_info"       // 获取完整的书籍信息

	API_SERVE_BOOK_FILE_BY_HASH string = "/api/book_library/serve_book_file_by_hash"
	API_SERVE_BOOK_FILE_BY_ID   string = "/api/book_library/serve_book_file_by_id"

	API_GET_BOOK_FIRST_PAGE_PDF string = "/api/book_library/get_book_first_page_pdf"
)

// apiServeBookFile

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
	ginserver.Listen(API_GET_BOOKS_INFO, apiGetBooksInfo, "返回某一页的书籍ID与书籍的名称、图像等属性。接受参数：page:int, page_size:int.")
	ginserver.Listen(API_GET_BOOKS_SHORT_INFO, apiGetBooksShortInfo, "返回某一页的书籍ID与书籍的名称、图像等属性。接受参数：page:int, page_size:int.")

	//apiServeBookFile
	ginserver.Listen(API_SERVE_BOOK_FILE_BY_HASH, apiServeBookFileByHash, "提供书籍文件的流式,接受book_file_hash作为输入")
	ginserver.Listen(API_SERVE_BOOK_FILE_BY_ID, apiServeBookFileById, "提供书籍文件的流式,接受book_id作为输入")

	// TEST!!!
	ginserver.Listen(API_GET_BOOK_FIRST_PAGE_PDF, apiGetBookFirstPagePdf, "获取书籍的首页PDF")
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
		BookId        string `json:"book_id" mapstructure:"book_id" gorm:"column:book_id"`
		BookName      string `json:"book_name" mapstructure:"book_name" gorm:"column:book_name"`
		BookImagePath string `json:"book_image_path" mapstructure:"book_imagepath" gorm:"column:book_image_path"`

		Author       string `json:"author" mapstructure:"author" gorm:"column:author"`
		Description  string `json:"description" mapstructure:"description" gorm:"column:description"`
		BookFileType string `json:"book_file_type" mapstructure:"book_file_type" gorm:"column:book_file_type"`
		BookFileHash string `json:"book_file_hash" mapstructure:"book_file_hash" gorm:"column:book_file_hash"`
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
		bookShortInfos[i].Author = book.Author
		bookShortInfos[i].Description = book.Description
		bookShortInfos[i].BookFileType = book.BookFileType
		bookShortInfos[i].BookFileHash = book.BookFileHash
	}

	// 返回这个数组
	c.JSON(http.StatusOK, gin.H{"books": bookShortInfos})
}

// 获取代表完整的书籍信息的数组
func apiGetBooksInfo(c *gin.Context) {
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
	idArr, err := GetAllBookIds(page, page_size) // String[], Error
	fmt.Println("  读取书籍ID:", idArr, err)
	if err != nil {
		// 如果出错，返回错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 根据已有ID数组，获取书籍信息。

	// 定义数组(切片)：
	bookInfos := make([]Book, len(idArr))

	// 遍历ID数组，获取书籍信息
	for i, book_id := range idArr {
		book, err := GetBookInfoById(book_id) // Book, Error
		if err != nil {
			// 如果出错，返回错误信息
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		bookInfos[i] = book
	}

	// 返回这个数组
	c.JSON(http.StatusOK, gin.H{"books": bookInfos})
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

// 提供书籍文件的流式传输服务, 接受 book_file_hash作为输入。
func apiServeBookFileByHash(c *gin.Context) {
	// 获取书籍文件
	// 书籍文件的哈希值
	bookFileHash := c.Query("book_file_hash")
	err := Bucket.ServeFile(c.Writer, c.Request, bookFileHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		color.Yellow("apiServeBookFile::ServeFile()::服务流失传输文件遇到了错误:", err)
		return
	}
}

// 提供书籍文件的流式传输服务, 接受 book_id作为输入。
func apiServeBookFileById(c *gin.Context) {
	// 获取书籍文件
	// 书籍ID
	bookID := c.Query("book_id")
	book, err := GetBookInfoById(bookID) // Book, Error
	if err != nil {
		// 如果出错，返回错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = Bucket.ServeFile(c.Writer, c.Request, book.BookFileHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		color.Yellow("apiServeBookFile::ServeFile()::服务流失传输文件遇到了错误: id=", bookID, ",错误：", err)
		return
	}
}

// 获取书籍的首页PDF
func apiGetBookFirstPagePdf(c *gin.Context) {
	// 获取书籍信息
	bookID := c.Query("book_id")         // 获取书籍ID
	book, err := GetBookInfoById(bookID) // Book, Error
	fmt.Println("  读取书籍信息:", book, err)

	if err != nil {
		// 如果出错，返回错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"GetBookInfoById::error": err.Error()})
		return
	}

	// type FirstPageInfo struct {
	// 	gorm.Model
	// 	BookId        string `json:"book_id" gorm:"column:book_id"`
	// 	FileType      string `json:"file_type" gorm:"column:file_type"`
	// 	FirstPageHash string `json:"first_page_hash" gorm:"column:first_page_hash"`
	// }
	// 书籍的hash和书籍首页文件的hash并不一样，可以通过数据库查询
	// 从FirstPageDB中查询书籍首页文件的hash，通过book_id 查询 first_page_hash

	// 根据book_id查询first_page_hash
	var firstPageInfo FirstPageInfo
	FirstPageDB.Where("book_id = ?", bookID).First(&firstPageInfo)
	bookFirstPageHash := firstPageInfo.FirstPageHash

	color.Red("apiGetBookFirstPagePdf::尝试获取首页哈希：", bookFirstPageHash)

	err = FirstPageBucket.ServeFile(c.Writer, c.Request, bookFirstPageHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		color.Yellow("apiGetBookFirstPagePdf::ServeFile()::获取书籍首页时，传输文件遇到了错误:", err)
		return
	}
}
