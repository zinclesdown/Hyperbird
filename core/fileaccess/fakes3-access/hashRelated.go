package fakes3access

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/blake2b"
)

// 计算指定路径的文件的哈希值
func (f *FS3Bucket) ComputeHash(path string) (string, error) {
	cleanPath := filepath.Clean(path) // 清理路径

	file, err := os.Open(cleanPath) // 打开文件
	if err != nil {
		fmt.Printf("Error opening file %q: %v\n", cleanPath, err)
		return "", err
	}
	defer file.Close()

	// 创建哈希计算器
	var hasher hash.Hash
	switch f.HashMethod {
	case "blake2b":
		hasher, _ = blake2b.New256(nil)
	case "md5":
		hasher = md5.New()
	case "sha1":
		hasher = sha1.New()
	case "sha256":
		hasher = sha256.New()
	case "sha512":
		hasher = sha512.New()
	default:
		return "", fmt.Errorf("unsupported hash method: %s", f.HashMethod)
	}

	// 读取文件内容并更新哈希计算器
	if _, err := io.Copy(hasher, file); err != nil {
		fmt.Printf("Error reading file %q: %v\n", cleanPath, err)
		return "", err
	}

	// 返回哈希值
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
