package fakes3access

import (
	"fmt"
	"io"
	"os"
)

// 自动化测试,当运行go run . 时,会通过test的Hello调用
func Test() {
	fmt.Println("\033[33m开始测试虚拟存储桶\033[0m")

	os.RemoveAll("./tmp/test/bucket") // 先删除测试路径下的所有文件

	test_bucket_create()
	test_bucket_read()
	test_bucket_db_functionalities()

	fmt.Println("\033[32m[虚拟存储桶的单元测试完成]\033[0m")
}

func test_bucket_create() {
	f := &FS3Bucket{} // 创建一个我们刚才写的Bucket

	bucket, err := f.CreateBucket("测试桶🪣", "./tmp/test/bucket/测试桶🪣的保存位置📦", Blake2b, 32) // 使用FS3Bucket实例调用CreateBucket方法

	fmt.Println("  创建测试 Bucket:", bucket)
	fmt.Println("  创建测试 ERR:", err)
}

func test_bucket_read() {
	fmt.Println("开始测试读取bucket")

	// 读取刚才的bucket
	f := &FS3Bucket{}
	bucket, err := f.LoadBucket("./tmp/test/bucket/测试桶🪣的保存位置📦")

	fmt.Println("  读取测试 Bucket:", bucket)
	fmt.Println("  读取测试 ERR:", err)
	fmt.Println("完成测试读取bucket")
}

func test_bucket_db_functionalities() {
	fmt.Println("开始测试数据库相关功能") // 测试数据库相关功能

	f := &FS3Bucket{} // 创建一个我们刚才写的Bucket
	bucket, err := f.LoadBucket("./tmp/test/bucket/测试桶🪣的保存位置📦")
	if err != nil {
		fmt.Println("读取bucket时遇到错误:", err)
		return
	}

	fmt.Println("> 写入文件测试 :")

	// 把文件加入桶中
	_, err = bucket.SaveFileFromPath("/home/zincles/Projects/Hyperbird/tmp/hello.txt", false)
	if err != nil {
		fmt.Println("保存文件时遇到错误:", err)
		return
	}
	_, err = bucket.SaveFileFromPath("/home/zincles/Projects/Hyperbird/tmp/hello2.txt", false)
	if err != nil {
		fmt.Println("保存文件时遇到错误:", err)
		return
	}
	_, err = bucket.SaveFileFromPath("/home/zincles/Projects/Hyperbird/tmp/hello3.txt", false)
	if err != nil {
		fmt.Println("保存文件时遇到错误:", err)
		return
	}

	// 读取文件
	fmt.Println("> 读取文件测试 :")
	hashs, err := bucket.GetAllFileHash()
	if err != nil {
		fmt.Println("读取文件时遇到错误:", err)
		return
	}
	printArray(hashs)

	// 获取第一个文件的大小
	size, err := bucket.GetFileSize(hashs[0])
	if err != nil {
		fmt.Println("获取文件大小时遇到错误:", err)
		return
	}
	fmt.Println("第一个文件的大小:", size)

	// 打开第二个文件,读取里面的内容作为字符串,并打印
	reader, err := bucket.GetFileReader(hashs[1]) // 获取io.Reader
	if err != nil {
		fmt.Println("打开文件时遇到错误:", err)
		return
	}
	data, err := io.ReadAll(reader) // 读取文件内容
	if err != nil {
		fmt.Println("读取文件时遇到错误:", err)
		return
	}
	fmt.Println("- 第二个文件的内容:\n", string(data))

	// 最后尝试删除第一个文件.
	err = bucket.DeleteFile(hashs[0])
	if err != nil {
		fmt.Println("删除文件时遇到错误:", err)
		return
	} else {
		fmt.Println("删除文件成功")
	}

}

// 调试用函数, 打印数组
func printArray(arr []string) {
	fmt.Println("[")
	for _, a := range arr {
		fmt.Println("	", a)
	}
	fmt.Println("]")
}
