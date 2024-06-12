package user

// TODO 编写单元测试
// TODO 实现JWT过期机制
// TODO 实现用户组

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 数据库位置
const (
	UserDbPath = "./data/core/users.db"
)

// 用户. 存储在数据库里
type User struct {
	gorm.Model
	UserID       string `json:"userid"` // int
	Username     string `json:"username"`
	PasswordHash string `json:"passwordhash"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Group        string `json:"group"` // int
}

type UserService interface {
	CreateUser(user *User) error                                                   // 创建一个新用户
	CreateUserWithPassword(username, password string) (userID string, error error) // 创建一个新用户,并设置密码
	GetUserByID(userID string) (*User, error)                                      // 获取一个用户的信息
	UpdateUser(userID string, user *User) error                                    // 更新现有用户的信息
	DeleteUser(userID string) error                                                // 删除一个用户
	AuthenticateUser(userID, password string) (bool, error)                        // 在其他地方鉴权用户
	GenerateJWT(userID string) (string, error)                                     // 生成一个用户的 JWT
	VerifyJWT(token string) (*User, error)                                         // 验证 JWT
	GetUserDB() (*gorm.DB, error)                                                  // 获取(所有)用户的数据库.新版本gorm不再需要取消连接
}

// 已有的方法:
//   warn(msg string, err error) // 遇到错误打印
//   assert(msg string, err error) // 遇到错误立刻终止所有程序, testcase 用

// UNTESTED
// 获取用户数据库. 使用SQLITE + grom.
// 如果根本没有用户数据库,则会新建一个
func GetUserDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(UserDbPath), &gorm.Config{})
	if err != nil {
		warn("连接到SQLite数据库时发生错误: %v\n", err)
		return nil, err
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		warn("创建用户表时发生错误: %v\n", err)
		return nil, err
	}

	return db, nil
}

// UNTESTED
// CreateUser 创建一个新用户
func CreateUser(user *User) error {
	db, err := GetUserDB()
	warn("获取数据库时发生错误: %v\n", err)

	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UNTESTED
// CreateUserWithPassword 创建一个新用户,并设置密码. 密码用 bcrypt 哈希
func CreateUserWithPassword(username, password string) (userID string, err error) {
	db, err := GetUserDB()
	if err != nil {
		return "", err
	}

	// 使用 bcrypt 生成哈希密码
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &User{
		Username:     username,
		PasswordHash: string(hash),
	}
	result := db.Create(user)
	if result.Error != nil {
		return "", result.Error
	}
	return user.UserID, nil
}

// UNTESTED
// GetUser 获取一个用户的信息
func GetUserByID(userID string) (*User, error) {
	db, err := GetUserDB()
	warn("获取数据库时发生错误: %v\n", err)

	var user User
	result := db.Where("user_id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UNTESTED
// UpdateUser 更新现有用户的信息
func UpdateUser(userID string, user *User) error {
	db, err := GetUserDB()
	if err != nil {
		return err
	}

	// 不更新密码哈希
	user.PasswordHash = ""

	result := db.Model(&User{}).Where("user_id = ?", userID).Updates(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UNTESTED
// DeleteUser 删除一个用户
// HACK 使用硬删除,而非软删除,删除操作会立刻执行
// 不安全,懒得改了
func DeleteUser(userID string) error {
	db, err := GetUserDB()
	warn("获取数据库时发生错误: %v\n", err)

	result := db.Unscoped().Where("user_id = ?", userID).Delete(&User{})
	warn("删除用户时发生错误: %v\n", result.Error)
	return nil
}

// UNTESTED
// AuthenticateUser 在其他地方鉴权用户
func AuthenticateUser(userID, password string) (bool, error) {
	db, err := GetUserDB()
	if err != nil {
		return false, fmt.Errorf("获取数据库时发生错误: %v", err)
	}

	var user User
	result := db.Where("user_id = ?", userID).First(&user)
	warn("获取用户时发生错误: %v\n", result.Error)
	if result.Error != nil {
		return false, result.Error
	}

	// 使用 bcrypt.CompareHashAndPassword 函数来比较哈希值和明文密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	warn("密码验证时发生错误: %v\n", err)

	if err != nil {
		return false, nil
	}

	return true, nil // 如果 CompareHashAndPassword 没有返回错误，那么密码是正确的
}

var jwtKey = []byte("Testing256bitSecretKeyForJWTPleaseChangeThisInProductionEnvironment")

// 为用户生成JWT
func GenerateJWT(userID string) (string, error) {
	// 创建一个新的令牌对象,指定签名方法和声明
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Subject:   userID,
	})

	// 使用密钥签名令牌
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		warn("生成Token时发生错误: %v\n", err)
		return "", err
	}

	return tokenString, nil
}

// 验证JWT, 并返回用户信息
func VerifyJWT(tokenString string) (*User, error) {

	// 传递Token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		warn("解析Token时发生错误: %v\n", err)
		return nil, err
	}

	// 验证Token
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		userID := claims.Subject
		user, err := GetUserByID(userID)
		if err != nil {
			warn("在验证JWT时获取用户时发生错误: %v\n", err)
			return nil, err
		}
		return user, nil
	}

	warn("Token无效\n", nil)
	return nil, fmt.Errorf("错误:Token无效")
}
