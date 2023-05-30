package models

import (
	"schoolChat/database"
	"schoolChat/util"
	"time"
)

// User 用户模型
type User struct {
	ID         int64  `json:"id"`          // 用户ID
	Username   string `json:"username"`    // 用户名
	Nickname   string `json:"nickname"`    // 昵称
	Avatar     string `json:"avatar"`      // 头像
	Password   string `json:"password"`    // 密码
	Status     int    `json:"status"`      // 状态
	CreateTime int64  `json:"create_time"` // 创建时间
	UpdateTime int64  `json:"update_time"` // 更新时间
}

// UserLogin 用户登录参数
type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AddUser @Title 新增用户
func (user *User) AddUser() (userID int64, err error) {
	user.Status = 1
	user.CreateTime = time.Now().UnixMilli()
	user.UpdateTime = time.Now().UnixMilli()
	user.Password = util.MD5(user.Password)

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

// CheckUser 检查用户是否存在
func CheckUser(username string) (user User) {
	database.GetDB().Where("username = ?", username).First(&user)
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
