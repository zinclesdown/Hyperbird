package booklibrary

import "io"

// 用户API接口,供前端用户调用
type UserAPI interface {
	Login(username string, password string) (UserContext, error)             // 登录接口，返回用户上下文
	GetAllBookIds(ctx UserContext, page int, pageSize int) ([]string, error) // 获取所有书籍 ID
	GetBookInfo(ctx UserContext, bookid string) (Book, error)                // 获取书籍信息
	GetBookFile(ctx UserContext, bookid string) (FileType, io.Reader, error) // 获取书籍文件
}

// 管理员API接口,供前端管理员调用
type AdminAPI interface {
	Authenticate(jwtToken string) (UserContext, error) // 鉴权接口，返回用户上下文
	AddBook(ctx UserContext, book Book) error          // 添加书籍
	UpdateBookInfo(ctx UserContext, book Book) error   // 更新书籍信息
	DeleteBook(ctx UserContext, bookid string) error   // 删除书籍
}
