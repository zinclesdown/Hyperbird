package directaccess

import (
	"os"
	"path/filepath"
)

// TODO
// 这是直接访问文件系统/网络系统的方法
func ServeFile(path string) {
}

// TEST
// 读取文件内容并直接返回字符串.
// 接受文件路径作为输入, 返回直接以文件形式读取该文件所获得的字符串.
func GetFileAsString(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	str := string(content) // 将文件内容转换为字符串
	return str, nil
}

// TEST
// 读取文件内容并直接返回字节切片.
// 接受文件路径作为输入, 返回直接以文件形式读取该文件所获得的字节切片.
func GetFileAsBytes(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// TESTME
// 获取文件夹下的所有文件和文件夹.
// 输入参数: 文件夹路径(string)
// 返回一个字符串数组, 包含文件和文件夹的绝对路径,
func GetFolderFiles(path string) ([]string, error) {
	// println("读取文件夹:", path)
	files, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(files)) // 存储文件和文件夹的绝对路径

	// 遍历文件和文件夹
	for _, file := range files {
		absPath, err := filepath.Abs(file)
		if err != nil {
			println("获取文件绝对路径失败:", err)
			continue
		}
		paths = append(paths, absPath) // 添加到路径数组
	}
	return paths, nil
}
