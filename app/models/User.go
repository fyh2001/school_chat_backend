package models

import (
	"schoolChat/database"
	"time"
)

// User 用户模型
type User struct {
	ID         int64  `json:"id"`          // 用户ID
	Email      string `json:"email"`       // 邮箱
	Phone      string `json:"phone"`       // 手机号
	Nickname   string `json:"nickname"`    // 昵称
	Gender     string `json:"gender"`      // 性别
	Avatar     string `json:"avatar"`      // 头像
	Background string `json:"background"`  // 个人主页背景图
	Status     int    `json:"status"`      // 状态
	CreateTime int64  `json:"create_time"` // 创建时间
	UpdateTime int64  `json:"update_time"` // 更新时间
}

type UserLoginByPhone struct {
	Phone string `json:"phone"` // 手机号
	Code  string `json:"code"`  // 验证码
}

// UserLoginByEmail 用户邮箱登录模型
type UserLoginByEmail struct {
	Email string `json:"email"` // 邮箱
	Code  string `json:"code"`  // 验证码
}

// UserRegisterByEmail 用户邮箱注册模型
type UserRegisterByEmail struct {
	Email    string `json:"email"`    // 邮箱
	Code     string `json:"code"`     // 验证码
	Nickname string `json:"nickname"` // 昵称
	Avatar   string `json:"avatar"`   // 头像
}

// AddUser @Title 新增用户
func (user *User) AddUser() (userID int64, err error) {
	user.Status = 1
	user.CreateTime = time.Now().UnixMilli()
	user.UpdateTime = time.Now().UnixMilli()

	result := database.GetDB().Create(&user) //  这里的DB变量是 database 包里定义的，Create 函数是 gorm包的创建数据API

	if result.Error != nil {
		err = result.Error
	}

	return userID, err // 返回新建数据的id 和 错误信息，在控制器里接收
}

// GetAllUser 获取全部用户信息
func GetAllUser() (users []User, err error) {
	result := database.GetDB().Find(&users)

	if result.Error != nil {
		err = result.Error
	}

	return users, err
}

// GetUserByUserId 根据用户ID获取用户信息
func GetUserByUserId(id any) (user User, err error) {
	result := database.GetDB().Where("id = ?", id).First(&user)

	if result.Error != nil {
		err = result.Error
	}

	return user, err
}

func GetUserByNickname(nickname string) (user User, err error) {
	result := database.GetDB().Where("nickname = ?", nickname).First(&user)

	if result.Error != nil {
		err = result.Error
	}

	return user, err
}

// GetUserByPhone 根据手机号获取用户信息
func GetUserByPhone(phone string) (user User, err error) {
	result := database.GetDB().Where("phone = ?", phone).First(&user)

	if result.Error != nil {
		err = result.Error
	}

	return user, err
}

// GetUserByEmail 根据邮箱获取用户信息
func GetUserByEmail(email string) (user User, err error) {
	result := database.GetDB().Where("email = ?", email).First(&user)

	if result.Error != nil {
		err = result.Error
	}

	return user, err
}

// CheckUser 检查用户是否存在
func CheckUser(email string) (user User) {
	database.GetDB().Where("email = ?", email).First(&user)
	return user
}

// UpdateUser 更新用户信息
func UpdateUser(id int64, user map[string]interface{}) (err error) {
	result := database.GetDB().Table("user").Where("id = ?", id).Updates(user)

	if result.Error != nil {
		err = result.Error
	}

	return err
}

func UpdateNickname(id int64, nickname string) (err error) {
	result := database.GetDB().Table("user").Where("id = ?", id).Update("nickname", nickname).Update("update_time", time.Now().UnixMilli())

	if result.Error != nil {
		err = result.Error
	}

	return err
}
