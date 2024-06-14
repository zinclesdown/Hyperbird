// TODO 施工中

package booklibrary

import (
	"fmt"
	FS3 "hyperbird/core/fileaccess/fakes3-access"
	"io"

	"gorm.io/driver/sqlite"
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

var Bucket *FS3.FS3Bucket // 书籍库的文件系统

var LibraryDB *gorm.DB // 书籍库的数据库

// 书籍文件的信息,保存在 booklibrary.db 数据库中
type Book struct {
	gorm.Model
	BookId          string `json:"bookid"`
	BookName        string `json:"bookname"`
	Author          string `json:"author"`
	Description     string `json:"description"`
	BookFileType    string `json:"bookfiletype"`
	BookFileHash    string `json:"bookfilehash"`
	AvailableGroups string `json:"availablegroups"` // 分割符为逗号的字符串,空格被视为字符的一部分
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

	AddBook(ctx UserContext, book Book) error // 添加书籍。书籍信息（Hash）务必准确，否则会导致文件无法访问而报错

	UpdateBookInfo(ctx UserContext, book Book) error
	DeleteBook(ctx UserContext, bookid string) error
}

// 初始化书籍库
// 全局变量Bucket会在这里被设置.
// 全局变量Bucket会在这里被设置.
func RunServer() {
	// 初始化书籍文件库
	f := &FS3.FS3Bucket{}
	if !f.HasBucket(BookLibraryBucketPath) { // 如果没有书籍文件库,则创建一个
		f.CreateBucket("booklibrary", BookLibraryBucketPath, FS3.Blake2b, 32)
	}
	bucket, err := f.LoadBucket(BookLibraryBucketPath) // 加载书籍库
	if err != nil {
		fmt.Println("初始化书籍库时遇到错误:", err)
		return
	}
	fmt.Println("书籍文件库初始化成功", bucket)
	Bucket = bucket

	// 初始化书籍数据库, 保存到全局变量里.
	LibraryDB, err = gorm.Open(sqlite.Open(BookLibraryDatabasePath), &gorm.Config{})
	if err != nil {
		fmt.Println("初始化书籍数据库时遇到错误:", err)
		return
	}

	err = LibraryDB.AutoMigrate(&Book{}) // 自动创建表
	if err != nil {
		fmt.Println("自动迁移数据库时遇到错误:", err)
		return
	}
}

// 获取库中所有电子书籍的ID
func (b *Book) GetAllBookIds(page int, pageSize int) ([]string, error) {
	var bookids []string
	err := LibraryDB.Model(&Book{}).Limit(pageSize).Offset(page*pageSize).Pluck("BookId", &bookids).Error

	if err != nil {
		return nil, err
	}
	return bookids, nil
}

// 向书籍库中添加一本书。book结构必须完整，文件hash必须正确。
func (b *Book) AddBook(book Book) error {

	// 检查书籍的哈希是否在FS3文件系统中存在。使用Bucket.HasFile()检查是否存在文件
	// 如果不存在,则返回错误
	if !Bucket.HasFile(book.BookFileHash) {
		fmt.Println("AddBook:书籍文件不存在")
		return fmt.Errorf("AddBook:书籍文件不存在")
	}

	err := LibraryDB.Create(&book).Error // 添加书籍
	if err != nil {
		return err
	}

	return nil
}
