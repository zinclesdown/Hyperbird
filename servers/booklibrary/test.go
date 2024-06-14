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

	color.Green("[图书管理系统测试完毕]")

}
