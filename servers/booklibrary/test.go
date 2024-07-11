package booklibrary

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// 运行测试前需执行的东西。例如清理测试文件夹
func PreTestBeforeServerStart() {
	// 删除测试文件夹 ./data/servers/booklibrary/books
	// 删除测试文件夹 ./data/servers/booklibrary/booklibrary.db

	fmt.Println("  删除测试文件夹 ./data/servers/booklibrary/books")
	os.RemoveAll("./data/servers/booklibrary/books")

	fmt.Println("  删除测试文件夹 ./data/servers/booklibrary/firstpage")
	os.RemoveAll("./data/servers/booklibrary/firstpage")

	fmt.Println("  删除测试文件夹 ./data/servers/booklibrary/booklibrary.db")
	os.RemoveAll("./data/servers/booklibrary/booklibrary.db")

	fmt.Println("  删除测试文件夹 ./data/servers/booklibrary/firstpage.db")
	os.RemoveAll("./data/servers/booklibrary/firstpage.db")

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
	// AddBook(Book{BookId: "TestBooksID2",
	// 	BookName:     "第二本测试书籍",
	// 	BookFileHash: file2.Hash})

	// 添加重复书籍
	// AddBook(Book{BookId: "TestBooksID3withfile2",
	// 	BookName:     "第三本克隆测试书籍",
	// 	BookFileHash: file2.Hash})

	fmt.Println(file2)

	// 读取书籍列表
	bookids, err = GetAllBookIds(0, 10)
	fmt.Println("  读取书籍ID:", bookids, err)

	// 读取书籍(1)
	// getTestBook, err := GetBookInfoById("TestBooksID")
	// fmt.Println("  读取书籍信息:", getTestBook)
	// fmt.Println("书籍名称：", getTestBook.BookName)
	// fmt.Println("书籍Hash：", getTestBook.BookFileHash)
	assert("  读取书籍信息遇到错误:", err)

	color.Green("[图书管理系统测试完毕]")

	AddCustomTestBooks()
}

// 添加自定义的测试用例图书。由于DMCA，这里不提供测试用例文件。你可以在测试用例文件夹内自行定义。
func AddCustomTestBooks() {
	color.Blue("  添加自定义的测试用例图书...")
	// 读取YAML(如果有)

	// 检查是否有yaml。没有则返回
	if _, err := os.Stat("./tests/booklibrary/custom/custom_testcases.yaml"); os.IsNotExist(err) {
		color.Yellow("  未找到自定义测试用例文件。")
		return
	}

	type customBookImportConfig struct {
		BookId       string `mapstructure:"book_id"`
		BookName     string `mapstructure:"book_name"`
		Author       string `mapstructure:"author"`
		Description  string `mapstructure:"description"`
		BookFilePath string `mapstructure:"book_file_path"`
	}

	type customImportConfig struct {
		Books []customBookImportConfig `mapstructure:"books"`
	}

	// 读取YAML
	// viper.SetConfigName("custom_testcases") // 配置文件名称（无扩展名）
	viper.SetConfigType("yaml")
	viper.SetConfigFile("./tests/booklibrary/custom/custom_testcases.yaml")

	if err := viper.ReadInConfig(); err != nil {
		color.Red("Error reading config file, %s", err)
		return

	}

	var config customImportConfig
	err := viper.Unmarshal(&config)
	if err != nil {
		color.Red("Unable to decode into struct, %s", err)
		return
	}

	// 输出读取的配置信息
	for _, book := range config.Books {
		color.Green("Book Name: %s, Book ID: %s, FileHash: %s\n", book.BookName, book.BookId, book.BookFilePath)

		// 添加书籍到桶里
		file, err := Bucket.SaveFileFromPath(book.BookFilePath, false)
		warn("添加书籍到桶里遇到了错误：", err)
		fmt.Println("向桶里添加了文件：", file.Hash)

		AddBook(Book{BookId: book.BookId,
			BookName:     book.BookName,
			Description:  book.Description,
			Author:       book.Author,
			BookFileHash: file.Hash})

	}
}
