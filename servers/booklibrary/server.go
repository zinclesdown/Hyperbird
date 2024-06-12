// TODO 施工中

package booklibrary

import (
	"fmt"
	FS3 "hyperbird/core/fileaccess/fakes3-access"
	"io"

	"gorm.io/gorm"
)

type FileType string

// 所有书籍存储在 ./data/servers/booklibrary/books/ 目录下的FS3文件系统中
// 配置数据库为 ./data/servers/booklibrary/booklibrary.db

const BookLibraryBucketPath = "./data/servers/booklibrary/books/"
const BookLibraryDatabasePath = "./data/servers/booklibrary/booklibrary.db"

// 可能的所有文件的类型
const (
	PDF  FileType = "pdf"
	EPUB FileType = "epub"
	TXT  FileType = "txt"
)

// 书籍文件的信息
type Book struct {
	gorm.Model
	BookId          string   `json:"bookid"`
	BookName        string   `json:"bookname"`
	Author          string   `json:"author"`
	Description     string   `json:"description"`
	BookFileType    FileType `json:"bookfiletype"`
	BookFileHash    string   `json:"bookfilehash"`
	AvailableGroups []string `json:"availablegroups"`
}

// 书籍库的访问接口
type BookLibraryServerAccessor interface {
	GetAllBookIds(page int, pageSize int) ([]string, error)
	GetBookInfo(bookid string) (Book, error)
	GetBookFileIOReader(bookid string) (FileType, io.Reader, error)
}

// 用户鉴权用的上下文
type UserContext struct {
	JwtToken string
}

// 书籍库的管理接口
type BookLibraryManageFunctions interface {
	Authenticate(jwtToken string) (UserContext, error)
	AddBook(ctx UserContext, book Book) error
	UpdateBookInfo(ctx UserContext, book Book) error
	DeleteBook(ctx UserContext, bookid string) error
}

// 初始化书籍库
func RunServer() {
	bucket, err := initBookLibraryDatabase() // 初始化数据库
	if err != nil {
		fmt.Println("初始化书籍库时遇到错误:", err)
		return
	}
	fmt.Println("书籍库初始化成功", bucket)
}

// 初始化数据库
func initBookLibraryDatabase() (*FS3.FS3Bucket, error) {
	f := &FS3.FS3Bucket{}

	if !f.HasBucket(BookLibraryBucketPath) { // 如果没有书籍库,则创建一个
		f.CreateBucket("booklibrary", BookLibraryBucketPath, FS3.Blake2b, 32)
	}

	bucket, err := f.LoadBucket(BookLibraryBucketPath) // 加载书籍库
	if err != nil {
		fmt.Println("初始化书籍库时遇到错误:", err)
		return bucket, err
	}
	return bucket, err
}
