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
	bookids, err := GetAllBookIds(1, 10)
	fmt.Println("  读取书籍ID:", bookids, err)

	// 添加测试用书籍（PDF）
	// color.White("  添加测试用书籍（PDF）...")
	pdf_path := "./tests/booklibrary/files/testbook.pdf"
	pdf_path2 := "tests/booklibrary/files/test2.pdf"

	// 添加书籍到桶里
	file, err := Bucket.SaveFileFromPath(pdf_path, false)
	warn("添加书籍到桶里遇到了错误：", err)
	fmt.Println("向桶里添加了文件：", file.Hash)
	AddBook(Book{BookId: "TestBooksID",
		BookName:     "第一本测试书籍",
		BookFileHash: file.Hash})

	file2, err := Bucket.SaveFileFromPath(pdf_path2, false)
	warn("添加书籍到桶里遇到了错误：", err)
	fmt.Println("向桶里添加了文件：", file.Hash)
	AddBook(Book{BookId: "TestBooksID2",
		BookName:     "第二本测试书籍",
		BookFileHash: file2.Hash})

	// 添加重复书籍
	AddBook(Book{BookId: "TestBooksID3withfile2",
		BookName:     "第三本克隆测试书籍",
		BookFileHash: file2.Hash})

	AddBook(Book{BookId: "TestBooksID4withfile2",
		BookName:     "第四本克隆测试书籍",
		BookFileHash: file2.Hash})
	AddBook(Book{BookId: "TestBooksID5withfile2",
		BookName:     "第五本克隆测试书籍",
		BookFileHash: file2.Hash})

	AddBook(Book{BookId: "TestBooksID6withfile2",
		BookName:     "第六本克隆测试书籍",
		BookFileHash: file2.Hash})

	AddBook(Book{BookId: "TestBooksID7withfile2",
		BookName:     "第七本克隆测试书籍",
		BookFileHash: file2.Hash})

	// 读取书籍列表
	bookids, err = GetAllBookIds(0, 10)
	fmt.Println("  读取书籍ID:", bookids, err)

	// 读取书籍(1)
	getTestBook, err := GetBookInfoById("TestBooksID")
	fmt.Println("  读取书籍信息:", getTestBook)
	fmt.Println("书籍名称：", getTestBook.BookName)
	fmt.Println("书籍Hash：", getTestBook.BookFileHash)
	assert("  读取书籍信息遇到错误:", err)

	color.Green("[图书管理系统测试完毕]")

}
