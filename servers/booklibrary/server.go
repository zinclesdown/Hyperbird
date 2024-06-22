// TODO 施工中

package booklibrary

import (
	"fmt"
	FS3 "hyperbird/core/fileaccess/fakes3-access"
	"io"
	"os"

	"github.com/fatih/color"
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

	// 基本信息
	BookId          string `json:"book_id" gorm:"column:book_id"`
	BookName        string `json:"book_name" gorm:"column:book_name"`
	BookImagePath   string `json:"book_imagepath" gorm:"column:book_image_path"`
	Author          string `json:"author" gorm:"column:author"`
	Description     string `json:"description" gorm:"column:description"`
	BookFileType    string `json:"book_file_type" gorm:"column:book_file_type"`
	BookFileHash    string `json:"book_file_hash" gorm:"column:book_file_hash"`
	AvailableGroups string `json:"available_groups" gorm:"column:available_groups"`

	// // 预览文件信息 (如果有的话), 单页PDF或者图片
	// PreviewFileType string `json:"preview_file_type" gorm:"column:preview_file_type"`
	// PreviewFileHash string `json:"preview_file_hash" gorm:"column:preview_file_hash"`
}

// 书籍库的访问接口
type BookLibraryServerAccessor interface {
	GetAllBookIds(page int, pageSize int) ([]string, error)
	GetBookInfoById(bookid string) (Book, error)
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
func InitServer() {
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

	InitServerFirstPage()
}

// 获取库中所有电子书籍的ID
func GetAllBookIds(page int, pageSize int) ([]string, error) {
	var bookids []string
	err := LibraryDB.Model(&Book{}).Limit(pageSize).Offset(page*pageSize).Pluck("BookId", &bookids).Error

	if err != nil {
		return nil, err
	}
	return bookids, nil
}

// 向书籍库中添加一本书。book结构必须完整，文件hash必须正确。
func AddBook(book Book) error {
	// 检查书籍的哈希是否在FS3文件系统中存在。使用Bucket.HasFile()检查是否存在文件。如果不存在,则返回错误
	if !Bucket.HasFile(book.BookFileHash) {
		fmt.Println("AddBook:书籍文件不存在")
		return fmt.Errorf("AddBook:书籍文件不存在")
	}

	err := LibraryDB.Create(&book).Error // 添加书籍
	if err != nil {
		return err
	}

	// 每当添加完书本时，向FirstPageBucket中添加书本的第一页。
	// 使用pdfcpu库将书本的第一页提取出来，检查哈希是否在预览桶中存在，
	// 如果不存在则添加到预览桶中，然后将预览哈希添加到书本的预览哈希字段中。
	// 如果存在，则直接将预览哈希添加到书本的预览哈希字段中。

	// FirstPageBucket *FS3.FS3Bucket // 书籍首页文件的桶
	// Bucket 		   *FS3.FS3Bucket // 书籍库的文件系统
	// LibraryDB 	   *gorm.DB // 书籍库的数据库
	// PdfFirstPageTmpPath // 临时目录位于 ./tmp/pdffirstpage/

	// 获取PDF文件的路径
	pdfpath, err := Bucket.GetFilePathReadOnly(book.BookFileHash) // *os.File
	color.Green("PDF文件路径:")
	fmt.Println(pdfpath)
	if err != nil {
		color.Red("AddBook:获取PDF文件路径时遇到错误")
		return err
	}

	// 打印FirstPageBucket的信息
	color.Green("FirstPageBucket:")
	fmt.Println(FirstPageBucket)

	// 提取第一页
	outputPos, err := ExtractFirstPageWithPdfCpuFile(pdfpath, PdfFirstPageTmpPath)
	if err != nil {
		color.Red("AddBook:提取第一页时遇到错误")
		return err
	}
	println("第一页文件路径:", outputPos)

	color.Yellow("尝试", outputPos)
	fileinfo, err := FirstPageBucket.SaveFileFromPath(outputPos, false)
	if err != nil {
		color.Red("AddBook:上传第一页文件时遇到错误")
		return err
	}
	fmt.Println("第一页文件上传成功:", fileinfo)

	// 接下来，还需要在FirstPageDB中添加一条记录
	// type FirstPageInfo struct {
	// 	gorm.Model
	// 	BookId        string `json:"book_id" gorm:"column:book_id"`
	// 	FileType      string `json:"file_type" gorm:"column:file_type"`
	// 	FirstPageHash string `json:"first_page_hash" gorm:"column:first_page_hash"`
	// }

	// 书籍的ID
	bookid := book.BookId
	// 书籍的第一页文件的哈希
	firstpagehash := fileinfo.Hash

	// 添加到FirstPageDB中
	err = FirstPageDB.Create(&FirstPageInfo{
		BookId:        bookid,
		FileType:      "pdf",
		FirstPageHash: firstpagehash,
	}).Error

	if err != nil {
		color.Red("AddBook:添加书籍首页信息时遇到错误")
		return err
	}

	color.Green("AddBook:书籍首页信息添加到数据库成功:")
	fmt.Println("  书籍ID:", bookid)
	fmt.Println("  书籍首页文件哈希:", firstpagehash)

	return nil
}

// 获取书籍信息
func GetBookInfoById(bookid string) (Book, error) {
	var book Book
	err := LibraryDB.Where("book_id = ?", bookid).First(&book).Error
	if err != nil {
		return Book{}, err
	}
	return book, nil
}

// 获取书籍文件的 *os.File
func GetBookFileIOReader(bookid string) (*os.File, error) {
	bucket := &FS3.FS3Bucket{}
	bucket, err := bucket.LoadBucket(BookLibraryBucketPath)

	if err != nil {
		color.Red("GetBookFile:加载书籍库时遇到错误")
		return nil, err
	}

	book, err := GetBookInfoById(bookid)
	if err != nil {
		color.Red("GetBookFile:获取书籍信息时遇到错误")
		return nil, err
	}

	file, err := bucket.OpenFile(book.BookFileHash)
	if err != nil {
		color.Red("GetBookFile:打开书籍文件时遇到错误")
		return nil, err
	}

	return file, nil
}

// 获取书籍文件的首页pdf的 *os.File
func GetBookFirstPageFileIOReader(bookid string) (*os.File, error) {

	// FirstPageDB的信息：
	// type FirstPageInfo struct {
	// 	gorm.Model
	// 	BookId        string `json:"book_id" gorm:"column:book_id"`
	// 	FileType      string `json:"file_type" gorm:"column:file_type"`
	// 	FirstPageHash string `json:"first_page_hash" gorm:"column:first_page_hash"`
	// }

	// 从该数据库获取对应的bookid的firstpagehash：
	firstpagehash := ""
	FirstPageDB.Where("book_id = ?", bookid).First(&firstpagehash)
	// 如果不存在，则返回错误
	if firstpagehash == "" {
		color.Red("GetBookFirstPageFile:获取书籍首页文件时遇到错误")
		return nil, fmt.Errorf("GetBookFirstPageFile:获取书籍首页文件时遇到错误")
	}

	// 从FirstPageBucket中获取对应的firstpagehash的文件：
	file, err := FirstPageBucket.OpenFile(firstpagehash)

	if err != nil {
		color.Red("GetBookFirstPageFile:打开书籍文件时遇到错误")
		return nil, err
	}

	return file, nil
}
