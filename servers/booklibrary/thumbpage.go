package booklibrary

import FS3 "hyperbird/core/fileaccess/fakes3-access"

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
