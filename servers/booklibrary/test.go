package booklibrary

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// 运行测试前需执行的东西。例如清理测试文件夹
func PreTestBeforeServerStart() {
	// 删除测试文件夹 ./data/servers/booklibrary/books
	// 删除测试文件夹 ./data/servers/booklibrary/booklibrary.db

	fmt.Println("  删除测试文件夹 ./data/servers/booklibrary/books")
	os.RemoveAll("./data/servers/booklibrary/books")
	fmt.Println("  删除测试文件夹 ./data/servers/booklibrary/booklibrary.db")
	os.RemoveAll("./data/servers/booklibrary/booklibrary.db")

}

func Test() {
	color.Blue("[开始测试图书管理系统]")

	// 查看桶内所有文件 (没有)
	hashs, err := Bucket.GetAllFileHash()
	fmt.Println("  首先读取测试 Bucket:", hashs, err)

	// 初始化书籍库
	fmt.Println("  初始化书籍库...")

	// 读取书籍
	b := &Book{}
	bookids, err := b.GetAllBookIds(1, 10)
	fmt.Println("  读取书籍ID:", bookids, err)

	// 添加测试用书籍（PDF）
	// color.White("  添加测试用书籍（PDF）...")
	pdf_path := "./tests/booklibrary/files/testbook.pdf"
	pdf_path2 := "tests/booklibrary/files/test2.pdf"

	// 添加书籍到桶里
	file, err := Bucket.SaveFileFromPath(pdf_path, false)
	warn("添加书籍到桶里遇到了错误：", err)
	fmt.Println("向桶里添加了文件：", file.Hash)
	b.AddBook(Book{BookId: "TestBooksID",
		BookFileHash: file.Hash})

	file2, err := Bucket.SaveFileFromPath(pdf_path2, false)
	warn("添加书籍到桶里遇到了错误：", err)
	fmt.Println("向桶里添加了文件：", file.Hash)
	b.AddBook(Book{BookId: "TestBooksID2",
		BookFileHash: file2.Hash})

	// 读取书籍列表
	bookids, err = b.GetAllBookIds(0, 10)
	fmt.Println("  读取书籍ID:", bookids, err)

	// 读取书籍(1)
	getTestBook, err := b.GetBookInfoById("TestBooksID")
	fmt.Println("  读取书籍信息:", getTestBook)
	fmt.Println("书籍名称：", getTestBook.BookName)
	fmt.Println("书籍Hash：", getTestBook.BookFileHash)
	assert("  读取书籍信息遇到错误:", err)

	color.Green("[图书管理系统测试完毕]")

}
