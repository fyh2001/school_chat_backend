package models

import (
	"schoolChat/database"
	"time"
)

// User 用户模型
type User struct {
	ID            int64  `json:"id"`                        // 用户ID
	Email         string `json:"email"`                     // 邮箱
	Phone         string `json:"phone"`                     // 手机号
	Nickname      string `json:"nickname"`                  // 昵称
	Gender        string `json:"gender"`                    // 性别
	Avatar        string `json:"avatar"`                    // 头像
	Signature     string `json:"signature"`                 // 个性签名
	Background    string `json:"background"`                // 个人主页背景图
	DiyBackground string `json:"diyBackground"`             // 个人主页自定义背景图
	Follows       int    `json:"follows" gorm:"default: 0"` // 关注数
	Fans          int    `json:"fans" gorm:"default: 0"`    // 粉丝数
	Likes         int    `json:"likes" gorm:"default: 0"`   // 获赞数
	Status        int    `json:"status"`                    // 状态
	CreateTime    int64  `json:"createTime"`                // 创建时间
	UpdateTime    int64  `json:"updateTime"`                // 更新时间
}

// LoginedUserResponse 登录用户响应模型
type LoginedUserResponse struct {
	ID            int64  `json:"id"`            // 用户ID
	Email         string `json:"email"`         // 邮箱
	Phone         string `json:"phone"`         // 手机号
	Nickname      string `json:"nickname"`      // 昵称
	Gender        string `json:"gender"`        // 性别
	Avatar        string `json:"avatar"`        // 头像
	Signature     string `json:"signature"`     // 个性签名
	Background    string `json:"background"`    // 个人主页背景图
	DiyBackground string `json:"diyBackground"` // 个人主页自定义背景图
	Follows       int    `json:"follows"`       // 关注数
	Fans          int    `json:"fans"`          // 粉丝数
	Likes         int    `json:"likes"`         // 获赞数
}

// UserResponse 用户响应模型
type UserResponse struct {
	ID            int64  `json:"id"`            // 用户ID
	Email         string `json:"email"`         // 邮箱
	Phone         string `json:"phone"`         // 手机号
	Nickname      string `json:"nickname"`      // 昵称
	Gender        string `json:"gender"`        // 性别
	Avatar        string `json:"avatar"`        // 头像
	Signature     string `json:"signature"`     // 个性签名
	Background    string `json:"background"`    // 个人主页背景图
	DiyBackground string `json:"diyBackground"` // 个人主页自定义背景图
	Follows       int    `json:"follows"`       // 关注数
	Fans          int    `json:"fans"`          // 粉丝数
	Likes         int    `json:"likes"`         // 获赞数
	IsFollow      bool   `json:"isFollow"`      // 是否关注
}

// UserLoginByPhone 用户手机登录模型
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

// UserFollows 用户关注模型
type UserFollows struct {
	ID       int64 `json:"id"`       // 主键ID
	UserID   int64 `json:"userId"`   // 用户ID
	FollowID int64 `json:"followId"` // 关注用户ID
}

// AddUser @Title 新增用户
func (user *User) AddUser() (userID int64, err error) {
	user.Status = 1
	user.CreateTime = time.Now().UnixMilli()
	user.UpdateTime = time.Now().UnixMilli()

	result := database.GetMySQL().Create(&user) //  这里的DB变量是 database 包里定义的，Create 函数是 gorm包的创建数据API

	if result.Error != nil {
		err = result.Error
	}
	userID = user.ID

	return userID, err // 返回新建数据的id 和 错误信息，在控制器里接收
}

// GetAllUser 获取全部用户信息
func GetAllUser() (users []User, err error) {
	result := database.GetMySQL().Find(&users)

	if result.Error != nil {
		err = result.Error
	}

	return users, err
}

// GetUserByLoginedUserId 根据登录用户ID获取用户信息
func GetUserByLoginedUserId(id int64) (user LoginedUserResponse, err error) {
	// 创建查询语句
	sql := `
		SELECT 
			user.id, 
			user.email,
			user.phone, 
			user.nickname, 
			user.gender, 
			user.avatar, 
			user.signature,
			user.background, 
			user.diy_background, 
			user.follows, 
			user.fans,
			user.likes
		FROM user
		WHERE user.id = ?
	`
	// 执行查询
	result := database.GetMySQL().Raw(sql, id).Scan(&user)

	if result.Error != nil {
		err = result.Error
	}

	return user, err
}

// GetUserByUserId 根据用户ID获取用户信息
func GetUserByUserId(LoginedUserId, userId int64) (user UserResponse, err error) {
	// 创建查询语句
	sql := `
		SELECT 
			user.id, 
			user.email,
			user.phone, 
			user.nickname, 
			user.gender, 
			user.avatar, 
			user.signature,
			user.background, 
			user.diy_background, 
			user.follows, 
			user.fans,
			user.likes,
			(CASE WHEN user_follows.follow_id IS NOT NULL THEN 1 ELSE 0 END) AS is_follow
		FROM user
		LEFT JOIN user_follows ON user_follows.user_id = ? AND user_follows.follow_id = user.id
		WHERE user.id = ?
`
	// 执行查询
	result := database.GetMySQL().Raw(sql, LoginedUserId, userId).Scan(&user)

	if result.Error != nil {
		err = result.Error
	}

	return user, err
}

// GetUserByNickname 根据昵称获取用户信息
func GetUserByNickname(nickname string) (user User, err error) {
	result := database.GetMySQL().Where("nickname = ?", nickname).First(&user)

	if result.Error != nil {
		err = result.Error
	}

	return user, err
}

// GetUserByPhone 根据手机号获取用户信息
func GetUserByPhone(phone string) (user User, err error) {
	result := database.GetMySQL().Where("phone = ?", phone).First(&user)

	if result.Error != nil {
		err = result.Error
	}

	return user, err
}

// GetUserByEmail 根据邮箱获取用户信息
func GetUserByEmail(email string) (user User, err error) {
	result := database.GetMySQL().Where("email = ?", email).First(&user)

	if result.Error != nil {
		err = result.Error
	}

	return user, err
}

// UpdateUser 更新用户信息
func UpdateUser(id int64, user map[string]interface{}) (err error) {
	result := database.GetMySQL().Table("user").Where("id = ?", id).Updates(user)

	if result.Error != nil {
		err = result.Error
	}

	return err
}

// UpdateNickname 更新用户昵称
func UpdateNickname(id int64, nickname string) (err error) {
	result := database.GetMySQL().Table("user").Where("id = ?", id).Update("nickname", nickname).Update("update_time", time.Now().UnixMilli())

	if result.Error != nil {
		err = result.Error
	}

	return err
}

// UpdateMail 更新用户邮箱
func UpdateMail(id int64, mail string) (err error) {
	result := database.GetMySQL().Table("user").Where("id = ?", id).Update("email", mail).Update("update_time", time.Now().UnixMilli())

	if result.Error != nil {
		err = result.Error
	}

	return err
}

// UpdatePhone 更新用户手机号
func UpdatePhone(id int64, phone string) (err error) {
	result := database.GetMySQL().Table("user").Where("id = ?", id).Update("phone", phone).Update("update_time", time.Now().UnixMilli())

	if result.Error != nil {
		err = result.Error
	}

	return err
}

// DeletePhone 解绑用户手机号
func DeletePhone(id int64) (err error) {
	result := database.GetMySQL().Table("user").Where("id = ?", id).Update("phone", "").Update("update_time", time.Now().UnixMilli())

	if result.Error != nil {
		err = result.Error
	}

	return err
}
