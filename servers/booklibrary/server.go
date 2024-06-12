package booklibrary

// TODO 使用GIN实现图书馆的API服务器

import (
	"github.com/gin-gonic/gin"
)

type BookLibraryServerAccessor interface {
	// 服务器访问器接口
}

func RunServer(apiUrl string, apiPort int) {
	// 服务器启动

}

func GetBooks(c *gin.Context) {
	// 获取所有书籍
}

func GetBook(c *gin.Context) {
	// 获取一本书
}

func PostBook(c *gin.Context) {
	// 添加一本书
}

func PutBook(c *gin.Context) {
	// 更新一本书
}

func DeleteBook(c *gin.Context) {
	// 删除一本书
}
