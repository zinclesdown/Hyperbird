package booklibrary

import (
	"fmt"
	FS3 "hyperbird/core/fileaccess/fakes3-access"
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// 将书籍PDF的封面转化为单页PDF,供前端展示
// 有点蠢，但最简单

// 存储第一页文件的桶
const FirstPageBucketPath = "./data/servers/booklibrary/firstpage/"

// 存储第一页数据的桶
var FirstPageBucket *FS3.FS3Bucket

func InitServerFirstPage() {
	// 初始化书籍首页文件库
	f := &FS3.FS3Bucket{}
	if !f.HasBucket(FirstPageBucketPath) { // 如果没有书籍文件库,则创建一个
		f.CreateBucket("bookfirstpage", FirstPageBucketPath, FS3.Blake2b, 32)
	}
}

// TEST 需要测试
// 将PDF文件的第一页提取到临时路径。
// 注意：每次尝试提取时，会清除临时路径下的文件。
// 因此提取完毕后建议立刻移动到指定路径，以免被覆盖。
func ExtractFirstPageWithPdfCpuFile(pdf_pos string, work_dir string) (out_pos string, err error) {
	// 移除临时路径下的所有文件
	if os.RemoveAll(work_dir) != nil {
		fmt.Printf("清除临时路径时出错: %v\n", err)
		return
	}

	// 创建临时路径
	if os.Mkdir(work_dir, os.ModePerm) != nil {
		fmt.Printf("创建临时路径时出错: %v\n", err)
		return
	}

	// 读取PDF文件
	ctx, err := api.ReadContextFile(pdf_pos)
	if err != nil {
		fmt.Printf("读取PDF文件时出错: %v\n", err)
		return
	}

	// 提取第一页
	reader, err := api.ExtractPage(ctx, 1)
	if err != nil {
		fmt.Printf("提取第一页时出错: %v\n", err)
		return
	}

	// 获取PDF文件的名称，不含路径
	pdf_name := filepath.Base(pdf_pos)

	// 保存到临时路径
	err = api.WritePage(reader, work_dir, pdf_name, 1)
	if err != nil {
		fmt.Printf("保存到临时路径时出错: %v\n", err)
		return
	}

	// 获取临时路径下第一个文件的路径。为了防止不兼容，直接列举临时路径下所有文件，并取第一个。
	files, err := os.ReadDir(work_dir)
	if err != nil {
		fmt.Printf("读取临时路径下的文件时出错: %v\n", err)
		return
	}

	// 返回第一个文件的路径
	return filepath.Join(work_dir, files[0].Name()), nil

}
