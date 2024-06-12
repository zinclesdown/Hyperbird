package user

import "gorm.io/gorm"

// 数据库位置
const (
	// 数据库位置
	UserDbPath = "/data/users.db"
)

// 用户. 存储在数据库里
type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Group    string `json:"group"`
}

// 用户组. 存储在数据库里
type UserGroup struct {
	gorm.Model
	GroupName string `json:"group_name"`
	GroupDesc string `json:"group_desc"`
}
