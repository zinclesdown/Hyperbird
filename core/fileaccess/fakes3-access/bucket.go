package fakes3access

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	helpers "hyperbird/core/helpers"
)

// TODO 代码中有大量的零散的数据库写入操作. 这问题很大,会导致性能严重下降甚至丢数据.
// HACK 另外,删除的实现也是有问题的. 以后再改吧.

// 在本地的虚拟的S3文件系统，通过两层文件夹(每层分别保存Hash的前两位，共4位)来模拟S3的Bucket和Object
// 文件夹分层，减少单个文件夹下的文件数量.该系统是为了模拟S3的Bucket和Object的文件系统
//
// 每个被视作FakeS3Folder的文件夹，都是一个Bucket，它包括以下结构：
// BucketFolder
// 	|- files
// 		|- AB
// 			|- CD
// 				|- ABCDE....DFGAS  // 文件名称为64个字符的Hash, 无扩展名
// 	|- bucket.json  // 元数据
//  |- filedb.db  // 文件数据库

// filedb.db的结构. 该db是对文件存储位置的缓存, 即使损坏或被删除也可以被重新构建.
type fileDB struct { // 存储于filedb.db中的文件的元数据
	gorm.Model
	Hash     string    // 文件的Hash，用于标识文件，使用32字符的哈希值
	Name     string    // 文件的名称，记录文件被更名为hash前的名称
	Path     string    // 文件的本地位置，用于提供文件服务
	Mime     string    // 文件的MIME类型，在添加到数据库时自动填充
	ExpireAt time.Time // 过期时间
}

// 虚拟S3的元数据.
type FS3Bucket struct {
	Directory  string     // 文件夹路径,文件夹包含fs3metadata.json
	BucketName string     // 容器名称
	HashMethod HashMethod // 计算Hash的方法 blake2b, md5, sha1, sha256, sha512
	HashLength int        // Hash的长度 32字符的哈希值，默认建议blake2b 32
	CreatedAt  time.Time  // 日期
}

// TODO 还没实现完 -- 甚至还没设计完! 这些接口不一定对的上下面实现的部分.
// 与FakeS3文件系统交互的接口. 供Bucket实现.
// 所有读写操作都是实时的.
// 我为什么要脑门一热写这么个设计?
// TODO 搞清楚自己到底在写什么
type FS3FileAccesser interface {

	// TEST 创建一个新的Bucket.
	CreateBucket(bucketName, directory string, hashMethod HashMethod, hashLength int) (*FS3Bucket, error)
	HasBucket(directory string) bool                                // TEST 检查指定名称的Bucket是否存在
	LoadBucket(directory string) (*FS3Bucket, error)                // TEST 加载指定名称的Bucket
	RecreateDB() error                                              // 重建数据库 重建DB,不会扫描,根据文件夹中的文件重新创建数据库. 也可以拿来初始化新桶
	RescanDB() error                                                // 重新扫描数据库  更新DB,不会删除文件,根据文件夹中的文件更新数据库
	SaveFileFromIO(data *os.File, filename string) (*fileDB, error) // 从IO将数据保存，并返回一个表示该文件的FileDB实例
	SaveFileFromPath(path string, cut bool) (*fileDB, error)        // 从指定路径将数据保存，并返回一个表示该文件的FileDB实例. cut为true时,使用move而非copy
	DeleteFile(hash string) error                                   // 删除指定哈希值的文件
	// SetExpire(hash string, expireAt time.Time) error         // 设置指定哈希值的文件的过期时间
	GetAllFileHash() (hashs []string, int64 error) // 返回所有的文件的哈希值
	GetFileSize(hash string) (int64, error)        // 返回指定哈希值的文件的大小
	// ClearExpired() error                                     // 删除所有过期的文件
	ComputeHash(path string) (string, error) // 计算指定路径的文件的哈希值,属性已经在FS3Bucket中定义
	GetFileDatabase() (*gorm.DB, error)      // 返回文件数据库
	HasFile(hash string) bool                // 检查指定哈希值的文件是否存在

	PrintBucketStatus() // 打印桶的状态

	OpenFile(hash string) (*os.File, error)                              // 返回指定哈希值的文件的数据
	ServeFile(w http.ResponseWriter, r *http.Request, hash string) error // 提供流式提供文件的功能，供前端的视频播放器/音频播放器/PDF查看器等使用

	GetFileNumber() (int, error) // 获取桶内文件的数量
}

// ========================================================================
//
//
//
// ===                      实现FS3FileAccesser接口                      ===
//
//
//
// ========================================================================

// 获取桶内文件的数量
func (f *FS3Bucket) GetFileNumber() (int, error) {
	db, err := f.GetFileDatabase()
	if err != nil {
		return 0, err
	}

	var files []fileDB
	if err := db.Find(&files).Error; err != nil {
		return 0, err
	}

	return len(files), nil
}

// UNTESTED!
// 从IO将数据保存，并返回一个表示该文件的FileDB实例
func (f *FS3Bucket) SaveFileFromIO(r *os.File, filename string) (*fileDB, error) {
	// 获取数据库连接
	db, err := f.GetFileDatabase()
	if err != nil {
		return nil, err
	}

	// 创建文件
	file, err := os.Create(f.Directory + "/" + filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 将数据从 io.Reader 写入到文件
	_, err = io.Copy(file, r)
	if err != nil {
		return nil, err
	}

	// 计算文件的哈希值
	hash, err := f.ComputeHash(f.Directory + "/" + filename)
	if err != nil {
		return nil, err
	}

	// 计算文件MIME
	mime, err := helpers.GetMimeFromFile(file)
	if err != nil {
		color.Red("计算文件MIME时发生错误: %v", err)
	}

	// 创建一个新的 fileDB 对象
	fileDB := &fileDB{
		Hash: hash,
		Path: f.Directory + "/" + filename,
		Name: filename,
		Mime: mime,
	}

	// 将 fileDB 对象保存到数据库
	if err := db.Create(fileDB).Error; err != nil {
		return nil, err
	}

	return fileDB, nil
}

// OpenFile 打开指定Hash的文件，返回*os.File
func (f *FS3Bucket) OpenFile(hash string) (*os.File, error) {
	db, err := f.GetFileDatabase()
	if err != nil {
		return nil, err
	}

	var file fileDB
	if err := db.Where("hash = ?", hash).First(&file).Error; err != nil {
		return nil, fmt.Errorf("在数据库中找不到哈希为 %s 的文件: %v", hash, err)
	}

	reader, err := os.Open(file.Path)
	if err != nil {
		return nil, fmt.Errorf("打开文件 %s 时发生错误: %v", file.Path, err)
	}

	return reader, nil
}

func PrintBucketStatus(f *FS3Bucket) {

	fileNum, err := f.GetFileNumber()
	if err != nil {
		fmt.Println("获取文件数量时发生错误:", err)
		return
	}

	color.Cyan("----- PrintBucketStatus: 桶的状态：-----")
	fmt.Println("  > BucketName:", f.BucketName)
	fmt.Println("  > Directory:", f.Directory)
	fmt.Println("  > HashMethod:", f.HashMethod)
	fmt.Println("  > HashLength:", f.HashLength)
	fmt.Println("  > CreatedAt:", f.CreatedAt)

	fmt.Println("  > 桶内文件数量：", fileNum)

	// 遍历每个桶内文件，输出filedb.db中的所有记录
	db, err := f.GetFileDatabase()
	if err != nil {
		fmt.Println("连接到数据库时发生错误:", err)
		return
	}

	var files []fileDB
	if err := db.Find(&files).Error; err != nil {
		fmt.Println("查询数据库时发生错误:", err)
		return
	}

	fmt.Println("  > 文件数量:", len(files))
	for _, file := range files {
		fmt.Printf(" - Name:[%s]  Mime:[%s]  Hash:[%s] \n", file.Name, file.Mime, file.Hash)
	}

	color.Cyan("---------------------------------------")
}

// 提供流式提供文件的功能，供前端的视频播放器/音频播放器/PDF查看器等使用
func (f *FS3Bucket) ServeFile(w http.ResponseWriter, r *http.Request, hash string) error {
	file, err := f.OpenFile(hash)
	if err != nil {
		return err
	}

	// 获取文件信息
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// 设置必要的响应头
	contentType, err := helpers.GetMimeFromFile(file)

	if err != nil {
		color.Red("ServeFile:获取文件MIME时发生错误: %v", err)
	}

	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(info.Name()))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}
	w.Header().Set("Content-Type", contentType)

	// 使用 http.ServeContent 发送文件
	http.ServeContent(w, r, info.Name(), info.ModTime(), file)

	return nil
}

func (f *FS3Bucket) GetAllFileHash() ([]string, error) {
	// 获取数据库连接
	db, err := f.GetFileDatabase()
	if err != nil {
		return nil, err
	}

	var files []fileDB
	if err := db.Find(&files).Error; err != nil {
		return nil, err
	}

	var hashes []string
	for _, file := range files {
		hashes = append(hashes, file.Hash)
	}

	return hashes, nil
}

// 获取一个桶内的文件的大小
func (f *FS3Bucket) GetFileSize(hash string) (int64, error) {
	// 获取数据库连接
	db, err := f.GetFileDatabase()
	if err != nil {
		return 0, err
	}

	var file fileDB
	if err := db.Where("hash = ?", hash).First(&file).Error; err != nil {
		return 0, err
	}

	info, err := os.Stat(file.Path)
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

// 获取一个桶的数据库
func (f *FS3Bucket) GetFileDatabase() (*gorm.DB, error) {
	// 连接到SQLite数据库
	db, err := gorm.Open(sqlite.Open(f.Directory+"/filedb.db"), &gorm.Config{})
	if err != nil {
		fmt.Printf("连接到SQLite数据库时发生错误: %v\n", err)
		return nil, err
	}

	return db, nil
}

// 实现CreateBucket接口.
// Directory必须要么不存在, 要么是空文件夹. 所有桶内的文件和桶索引数据库都会被存储在这个directory中.
func (f *FS3Bucket) CreateBucket(bucketName, directory string, hashMethod HashMethod, hashLength int) (*FS3Bucket, error) {
	_, err := os.Stat(directory) // 首先检查是否有目录

	// 判断
	switch {
	case err != nil && os.IsNotExist(err): // 目录不存在,则尝试创建目录. 如果创建失败则返回错误
		fmt.Println("目录不存在,创建目录:", directory)
		err := os.MkdirAll(directory, 0755)
		if err != nil {
			fmt.Println("创建目录时遇到错误:", err)
			return nil, fmt.Errorf("创建目录时遇到错误: %w", err)
		}
	case err != nil: // 其他错误
		fmt.Println("遇到了其他错误:", err)
		return nil, fmt.Errorf("遇到了其他错误: %w", err)
	default: // 目录存在
		fmt.Println("目录已存在")
	}

	// 创建bucket结构
	bucket := &FS3Bucket{
		Directory:  directory,  // string
		BucketName: bucketName, // string
		HashMethod: hashMethod, // string
		HashLength: hashLength, // int
		CreatedAt:  time.Now(), // time
	}

	// 把该结构体作为json直接保存到dir/bucket.json里,方便后续读取信息.
	bucketjson, err := os.Create(directory + "/bucket.json")
	if err != nil {
		fmt.Println("创建bucket.json时遇到错误:", err)
		return nil, fmt.Errorf("创建bucket.json时遇到错误: %w", err)
	}
	defer bucketjson.Close()

	// 将bucket结构转换为JSON字符串，然后写入文件
	bucketBytes, err := json.MarshalIndent(bucket, "", "  ")
	if err != nil {
		fmt.Println("格式化bucket结构供json写入时遇到错误:", err)
		return nil, fmt.Errorf("格式化bucket结构供json写入时遇到错误: %w", err)
	}

	// 写入bucket.json
	_, err = bucketjson.Write(bucketBytes)
	if err != nil {
		fmt.Println("写入bucket.json时遇到错误:", err)
		return nil, fmt.Errorf("写入bucket.json时遇到错误: %w", err)
	}

	// 初始化桶文件数据库
	bucket.RecreateDB()

	return bucket, nil
}

// TEST
// 实现HasBucket接口. 检查指定名称的Bucket是否存在.
func (f *FS3Bucket) HasBucket(directory string) bool {
	_, err := os.Stat(directory + "/bucket.json")
	return err == nil
}

// 实现LoadBucket接口. 加载现有Bucket.
func (f *FS3Bucket) LoadBucket(directory string) (*FS3Bucket, error) {
	// 读取bucket.json
	bucketBytes, err := os.ReadFile(directory + "/bucket.json")
	if err != nil {
		fmt.Println("读取bucket.json时遇到错误:", err) // 添加调试信息
		return nil, fmt.Errorf("读取bucket.json时遇到错误: %w", err)
	}

	// 将JSON字符串解码为FS3Bucket结构
	bucket := &FS3Bucket{}
	err = json.Unmarshal(bucketBytes, bucket)
	if err != nil {
		fmt.Println("解码bucket.json时遇到错误:", err) // 添加调试信息
		return nil, fmt.Errorf("解码bucket.json时遇到错误: %w", err)
	}

	fmt.Println("解码后的bucket:", bucket) // 添加调试信息
	return bucket, nil
}

// 重建或初始化桶的数据库
func (f *FS3Bucket) RecreateDB() error {
	// 连接到SQLite数据库
	db, err := f.GetFileDatabase()
	if err != nil {
		fmt.Printf("Error connecting to SQLite database: %v\n", err)
		return err
	}

	// 删除旧的fileDB表
	err = db.Migrator().DropTable(&fileDB{})
	if err != nil {
		fmt.Printf("Error dropping old fileDB table: %v\n", err)
		return err
	}

	// 创建新的fileDB表
	err = db.Migrator().CreateTable(&fileDB{})
	if err != nil {
		fmt.Printf("Error creating new fileDB table: %v\n", err)
		return err
	}

	// 黄色高亮输出
	fmt.Printf("\033[33m已重建filedb.db\033[0m\n")

	return nil
}

// 更新桶的数据库. 这将扫描桶中的所有文件，并将它们添加到数据库中.
func (f *FS3Bucket) RescanDB() error {
	// 连接到SQLite数据库
	db, err := f.GetFileDatabase()
	if err != nil {
		fmt.Printf("连接到SQLite数据库时发生错误: %v\n", err)
		return err
	}

	// 清理并检查路径
	cleanPath := filepath.Clean(f.Directory + "/files")
	if !strings.HasPrefix(cleanPath, f.Directory) {
		return fmt.Errorf("路径不在预期的目录中: %s", cleanPath)
	}

	// 遍历files文件夹
	err = filepath.Walk(cleanPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("访问路径 %q 时发生错误: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		// 计算文件的哈希值
		hash, err := f.ComputeHash(path)
		if err != nil {
			fmt.Printf("计算文件 %q 的哈希值时发生错误: %v\n", path, err)
			return err
		}

		// 检查数据库中是否已经存在具有相同哈希值的文件
		var file fileDB
		if err := db.Where("hash = ?", hash).First(&file).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				fmt.Printf("查询数据库时发生错误: %v\n", err)
				return err
			}
		} else {
			// 如果找到了具有相同哈希值的文件，跳过这个文件
			return nil
		}

		// 创建新的fileDB实例
		file = fileDB{
			Hash:     hash,
			Name:     filepath.Base(path),
			Path:     path,
			ExpireAt: time.Now().Add(24 * time.Hour), // 设置过期时间为24小时后
		}

		// 保存fileDB实例到数据库
		result := db.Create(&file)
		if result.Error != nil {
			fmt.Printf("将fileDB实例保存到数据库时发生错误: %v\n", result.Error)
			return result.Error
		}

		return nil
	})

	if err != nil {
		fmt.Printf("遍历路径 %q 时发生错误: %v\n", cleanPath, err)
		return err
	}

	return nil
}

// 使用直接路径,把文件加入桶中. 直接放在files/下, 向数据库加入一条记录.
// 从指定路径将数据保存，并返回一个表示该文件的FileDB实例. cut为true时,使用move而非copy
// SaveFileFromPath 将文件保存到桶中，并返回一个表示该文件的fileDB实例。
func (f *FS3Bucket) SaveFileFromPath(path string, cut bool) (*fileDB, error) {
	// 验证文件路径
	if strings.Contains(path, "..") {
		return nil, fmt.Errorf("文件路径包含非法字符: '..'")
	}

	// 连接到SQLite数据库
	db, err := f.GetFileDatabase()
	if err != nil {
		return nil, fmt.Errorf("连接到SQLite数据库时发生错误: %v", err)
	}

	// 计算文件的哈希值
	hash, err := f.ComputeHash(path)
	if err != nil {
		return nil, fmt.Errorf("计算文件 %q 的哈希值时发生错误: %v", path, err)
	}

	// 计算文件的MIME
	mime, err := helpers.GetMimeFromPath(path)
	if err != nil {
		color.Red("计算文件MIME时发生错误: %v", err)
	}

	// 检查数据库中是否已经存在具有相同哈希值的文件
	var file fileDB
	if err := db.Where("hash = ?", hash).First(&file).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("查询数据库时发生错误: %v", err)
		}
	} else {
		// 如果找到了具有相同哈希值的文件，返回这个文件的fileDB实例
		return &file, nil
	}

	// 根据哈希值的前四位创建文件的存储路径
	destDir := filepath.Join(f.Directory, "files", hash[:2], hash[2:4])
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("创建文件存储路径时发生错误: %v", err)
	}

	destPath := filepath.Join(destDir, hash)

	// 将文件复制或移动到新的存储路径
	if cut {
		err = os.Rename(path, destPath)
		if err != nil {
			return nil, fmt.Errorf("移动文件时发生错误: %v", err)
		}
	} else {
		srcFile, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("打开源文件时发生错误: %v", err)
		}
		defer srcFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return nil, fmt.Errorf("创建目标文件时发生错误: %v", err)
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return nil, fmt.Errorf("复制文件时发生错误: %v", err)
		}
	}

	// 创建新的fileDB实例
	file = fileDB{
		Hash: hash,
		Name: filepath.Base(path),
		Path: destPath,
		Mime: mime,
		// ExpireAt: time.Now().Add(24 * time.Hour), // 设置过期时间为24小时后
	}

	// 保存fileDB实例到数据库
	result := db.Create(&file)
	if result.Error != nil {
		return nil, fmt.Errorf("将fileDB实例保存到数据库时发生错误: %v", result.Error)
	}

	return &file, nil
}

// HACK 注意:使用了Unscoped()方法,这意味着删除的文件将永远不会被恢复.
// 我还没想好具体怎么设计, 但这个方法是有问题的.
// 长远来看,这对数据库的性能和稳定性都是有问题的!
// 以后再改吧.
func (f *FS3Bucket) DeleteFile(hash string) error {
	// 连接到SQLite数据库
	db, err := f.GetFileDatabase()
	if err != nil {
		fmt.Printf("连接到SQLite数据库时发生错误: %v\n", err)
		return err
	}

	var file fileDB
	if err := db.Where("hash = ?", hash).First(&file).Error; err != nil {
		fmt.Printf("在数据库中找不到哈希为 %s 的文件: %v\n", hash, err)
		return err
	}

	// 从数据库中完全删除记录
	if err := db.Unscoped().Delete(&file).Error; err != nil {
		return err
	}

	if err := os.Remove(file.Path); err != nil {
		return err
	}

	return nil
}

// TEST
func (f *FS3Bucket) HasFile(hash string) bool {
	// 连接到SQLite数据库
	db, err := f.GetFileDatabase()
	if err != nil {
		fmt.Printf("连接到SQLite数据库时发生错误: %v\n", err)
		return false
	}

	var fileRecord fileDB
	if err := db.Where("hash = ?", hash).First(&fileRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		fmt.Printf("查询数据库时发生错误: %v\n", err)
		return false
	}

	// 确认是否能访问
	if _, err := os.Stat(fileRecord.Path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		fmt.Printf("检查文件 %q 时发生错误: %v\n", fileRecord.Path, err)

		return false
	}

	return true
}
