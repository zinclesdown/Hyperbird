package helpers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/fatih/color"
)

// 根据文件，获取文件的MIME类型
func GetMimeFromFile(file *os.File) (string, error) {
	// 读取前512字节
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// 获取MIME类型
	contentType := http.DetectContentType(buffer[:n])
	fmt.Println("Content Type:", contentType)
	return contentType, nil
}

// 根据文件路径，获取文件的MIME类型
func GetMimeFromPath(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			// 这里可以使用你的日志库来记录错误
			color.Red("GetMimeFromPath:关闭文件遇到错误: ", closeErr)
		}
	}()

	return GetMimeFromFile(file)
}
