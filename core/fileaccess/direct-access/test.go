package directaccess

import "fmt"

func Test() {
	testDirectAccess()
}

func testDirectAccess() {
	println("测试文件直接访问...")
	_, err := GetFolderFiles("./tmp")

	if err != nil {
		fmt.Println("读取文件夹失败:", err)
		return
	}
	_, err = GetFileAsString("./tmp/hello.txt")

	if err != nil {
		fmt.Println("读取文件夹失败:", err)
		return
	}
	println("完成")
}
