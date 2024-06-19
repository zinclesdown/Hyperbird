package fakes3access

import (
	"bytes"
	"fmt"

	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/fatih/color"
)

const (
	file1Path = "./tests/fs3bucket/files/hello.txt"
	file2Path = "./tests/fs3bucket/files/hello2.txt"
	file3Path = "./tests/fs3bucket/files/hello3.txt"

	fileToServePath = "tests/fs3bucket/files/helloforserve.txt"

	bucketBasePath = "./tests/fs3bucket/tmp/"
	bucketPath     = "./tests/fs3bucket/tmp/测试用桶🪣/"
)

// 自动化测试,当运行go run . 时,会通过test的Hello调用
func Test() {
	fmt.Println("\033[33m开始测试虚拟存储桶\033[0m")

	os.RemoveAll(bucketBasePath) // 先删除测试路径下的所有文件

	test_bucket_create()
	test_bucket_read()
	test_bucket_db_functionalities()
	testServeFile()

	fmt.Println("\033[32m[虚拟存储桶的单元测试完成]\033[0m")
}

func test_bucket_create() {
	f := &FS3Bucket{} // 创建一个我们刚才写的Bucket

	bucket, err := f.CreateBucket("测试桶🪣名称", bucketPath, Blake2b, 32) // 使用FS3Bucket实例调用CreateBucket方法

	fmt.Println("  创建测试 Bucket:", bucket)
	checkErr("创建桶遇到了错误:", err)
}

func test_bucket_read() {
	fmt.Println("开始测试读取bucket")

	// 读取刚才的bucket
	f := &FS3Bucket{}
	bucket, err := f.LoadBucket(bucketPath)

	fmt.Println("  读取测试 Bucket:", bucket)
	checkErr("读取桶遇到了问题:", err)
	fmt.Println("完成测试读取bucket")
}

func test_bucket_db_functionalities() {

	fmt.Println("开始测试数据库相关功能") // 测试数据库相关功能

	f := &FS3Bucket{} // 创建一个我们刚才写的Bucket
	bucket, err := f.LoadBucket(bucketPath)
	checkErr("读取遇到错误", err)

	fmt.Println("> 写入文件测试 :")

	// 把文件加入桶中
	_, err = bucket.SaveFileFromPath(file1Path, false)
	checkErr("保存文件时遇到错误:", err)
	_, err = bucket.SaveFileFromPath(file2Path, false)
	checkErr("保存文件时遇到错误:", err)
	_, err = bucket.SaveFileFromPath(file3Path, false)
	checkErr("保存文件时遇到错误:", err)

	// 读取文件，打印哈希
	fmt.Println("> 读取文件测试 :")
	hashs, err := bucket.GetAllFileHash()
	checkErr("读取文件时遇到错误:", err)
	printArray(hashs)

	// 获取第一个文件的大小
	size, err := bucket.GetFileSize(hashs[0])

	checkErr("获取文件大小时遇到错误:", err)

	fmt.Println("第一个文件的大小:", size)

	// 打开第二个文件,读取里面的内容作为字符串,并打印
	file, err := bucket.OpenFile(hashs[1]) // 获取io.Reader
	checkErr("打开文件时遇到错误:", err)

	data, err := io.ReadAll(file) // 读取文件内容
	checkErr("读取文件时遇到错误:", err)

	fmt.Println("- 第二个文件的内容:\n", string(data))

	// 最后尝试删除第一个文件.
	err = bucket.DeleteFile(hashs[0])
	checkErr("删除文件时遇到错误:", err)

	// 打印所有文件的哈希
	fmt.Println("> 读取文件测试 :")
	hashs, err = bucket.GetAllFileHash()
	checkErr("读取文件时遇到错误:", err)
	printArray(hashs)

	// 最后，尝试Serve第2个文件
	fmt.Println("> 读取文件测试 :")

	// 从这里开始编写单元测试
}

// TODO 本单元测试有待被测试，有点蠢
// 由copilot编写
func testServeFile() {
	// 创建一个 FS3Bucket 实例和一个 HTTP 服务器
	bucket := &FS3Bucket{}
	bucket, err := bucket.LoadBucket(bucketPath)
	if err != nil {
		log.Fatal("读取桶时遇到错误:尝试读取路径：", bucketPath, "但遇到了错误：", err)
	}

	// 打印桶的状态
	fmt.Println("桶状态：")
	PrintBucketStatus(bucket)

	// 保存一个额外的文件到桶中
	filedb, err := bucket.SaveFileFromPath(fileToServePath, false)

	hash := filedb.Hash

	if err != nil {
		log.Fatal("testServeFile::保存文件时遇到错误:", err)
	}

	// 创建一个 HTTP 服务器来测试 ServeFile 方法
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := bucket.ServeFile(w, r, hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	// 使用 http.Get 函数模拟客户端请求
	resp, err := http.Get(server.URL)
	if err != nil {
		log.Fatal("发送 HTTP 请求时遇到错误:", err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应的状态码
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("期望的状态码是 %d 但是得到的是 %d", http.StatusOK, resp.StatusCode)
	}

	// 检查 HTTP 响应的 Content-Type 头部
	expectedContentType := "application/octet-stream"
	if contentType := resp.Header.Get("Content-Type"); contentType != expectedContentType {
		log.Fatalf("期望的 Content-Type 是 %s 但是得到的是 %s", expectedContentType, contentType)
	}

	// 检查 HTTP 响应的正文
	expectedBody, err := os.ReadFile(fileToServePath)
	if err != nil {
		log.Fatalf("读取文件时遇到错误: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("读取响应正文时遇到错误: %v", err)
	}
	if !bytes.Equal(body, expectedBody) {
		log.Fatalf("期望的响应正文是 \n%s\n,但是得到的是 \n%s", expectedBody, body)
	}

	color.Green("[文件服务测试通过！]\n")
}

// 调试用函数, 打印数组
func printArray(arr []string) {
	fmt.Println("[")
	for _, a := range arr {
		fmt.Println("	", a)
	}
	fmt.Println("]")
}

// 检查错误,如果有错误则打印错误信息并终止程序
func checkErr(msg string, err error) {
	if err != nil {
		// 遇到错误后终止程序
		fmt.Printf("\033[31m%s\033[0m", msg) // 使用红色文本打印msg
		log.Fatal(err)
	}
}
