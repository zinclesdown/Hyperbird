package fakes3access

import (
	"fmt"
	"os"
)

// 自动化测试,当运行go run . 时,会通过test的Hello调用
func Test() {
	fmt.Println("\033[33m开始测试虚拟存储桶\033[0m")

	// 先删除测试路径下的所有文件
	os.RemoveAll("./tmp/test/bucket")

	test_bucket_create()
	test_bucket_read()
	test_bucket_db_functionalities()

	fmt.Println("\033[32m[虚拟存储桶的单元测试完成]\033[0m")
}

func test_bucket_create() {
	// 创建一个我们刚才写的Bucket
	f := &FS3Bucket{}

	// 使用FS3Bucket实例调用CreateBucket方法
	bucket, err := f.CreateBucket("测试桶🪣", "./tmp/test/bucket/测试桶🪣的保存位置📦", Blake2b, 32)

	fmt.Println("  创建测试 Bucket:", bucket)
	fmt.Println("  创建测试 ERR:", err)
}

func test_bucket_read() {
	fmt.Println("开始测试读取bucket")
	// 测试读取bucket

	// 读取刚才的bucket

	f := &FS3Bucket{}
	bucket, err := f.LoadBucket("./tmp/test/bucket/测试桶🪣的保存位置📦")

	fmt.Println("  读取测试 Bucket:", bucket)
	fmt.Println("  读取测试 ERR:", err)
	fmt.Println("完成测试读取bucket")
}

func test_bucket_db_functionalities() {
	// 测试数据库相关功能
	fmt.Println("开始测试数据库相关功能")

	// 创建一个我们刚才写的Bucket
	f := &FS3Bucket{}
	bucket, err := f.LoadBucket("./tmp/test/bucket/测试桶🪣的保存位置📦")

	if err != nil {
		fmt.Println("读取bucket时遇到错误:", err)
		return
	}

	fmt.Println("  写入文件测试 :")

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

	// 最后尝试删除第一个文件.
	err = bucket.DeleteFile("9fbf4ee4ac272cf28e69a7bb624c01f94872733375bc4d599f3018fa35108925")
	if err != nil {
		fmt.Println("删除文件时遇到错误:", err)
		return
	} else {
		fmt.Println("删除文件成功")
	}

}
