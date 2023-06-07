package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"schoolChat/app/models"
	Results "schoolChat/app/result"
	"schoolChat/database"
	"schoolChat/util"
	"strconv"
	"time"
)

var avatar = []string{
	"girl_1.svg",
	"girl_2.svg",
	"girl_3.svg",
	"girl_4.svg",
	"girl_5.svg",
	"girl_6.svg",
	"girl_7.svg",
	"girl_8.svg",
	"girl_9.svg",
	"girl_10.svg",
	"girl_11.svg",
	"boy_1.svg",
	"boy_2.svg",
	"boy_3.svg",
	"boy_4.svg",
	"boy_5.svg",
}

// GetAllUser 获取所有用户
func GetAllUser(c *gin.Context) {
	var json models.User
	err := c.BindJSON(&json)

	if err != nil {
		fmt.Printf("数据库链接错误: %v", err)
	}

	results, err := models.GetAllUser()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  false,
			"message": "抱歉未找到相关信息",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  true,
		"data":    results,
		"message": "查询成功",
	})
}

// GetUserByUserId 根据用户ID获取用户信息
func GetUserByUserId(c *gin.Context) {
	userId, _ := c.Get("id")

	result, err := models.GetUserByUserId(userId.(int64))

	if err != nil {
		c.JSON(200, Results.Err.Fail("查询失败"))
		return
	}

	c.JSON(200, Results.Ok.Success(result))
}

// LoginOrRegisterByPhone 手机登录或注册
func LoginOrRegisterByPhone(c *gin.Context) {
	var phone models.UserLoginByPhone

	err := c.BindJSON(&phone)
	if err != nil {
		fmt.Printf("获取参数失败: %v", err)
		c.JSON(200, Results.ErrBind.Fail("获取参数失败"))
		return
	}

	// 验证验证码是否正确
	if phone.Code == "" {
		c.JSON(200, Results.Err.Fail("验证码不能为空"))
		return
	} else {
		// 从redis中获取验证码
		phoneCode, _ := database.GetRedis().Get(context.Background(), phone.Phone).Result()
		if phoneCode != phone.Code {
			c.JSON(200, Results.Err.Fail("验证码错误"))
			return
		}
	}

	// 验证手机号是否存在
	userData, _ := models.GetUserByPhone(phone.Phone)

	// 如果手机号不存在则注册
	if userData.ID == 0 {
		user := models.User{Phone: phone.Phone, Nickname: phone.Phone, Avatar: avatar[rand.Intn(len(avatar))]}
		userId, err := user.AddUser()
		if err != nil {
			c.JSON(200, Results.Err.Fail("注册失败"))
			return
		}

		// 创建token
		token, err := util.GenerateToken(userId, phone.Phone, "")
		if err != nil {
			c.JSON(200, Results.ErrToken.Fail("token生成失败"))
			return
		}

		data := make(map[string]interface{})
		data["token"] = token
		data["id"] = userId
		data["phone"] = phone.Phone
		data["nickname"] = phone.Phone
		data["avatar"] = "userData.Avatar"

		// 删除redis中的验证码
		database.GetRedis().Del(context.Background(), phone.Phone)

		c.JSON(200, Results.Ok.Success(data))
		return
	} else {
		// 如果手机号存在则登录
		// 创建token
		token, err := util.GenerateToken(userData.ID, userData.Phone, "")
		if err != nil {
			c.JSON(200, Results.ErrToken.Fail("token生成失败"))
			return
		}

		data := make(map[string]interface{})
		data["token"] = token
		data["id"] = userData.ID
		data["phone"] = userData.Phone
		data["nickname"] = userData.Nickname
		data["avatar"] = userData.Avatar

		// 删除redis中的验证码
		database.GetRedis().Del(context.Background(), phone.Phone)

		c.JSON(200, Results.Ok.Success(data))
		return
	}

}

// LoginOrRegisterByMail 邮箱登录或注册
func LoginOrRegisterByMail(c *gin.Context) {
	var mail models.UserLoginByEmail

	err := c.BindJSON(&mail)
	if err != nil {
		fmt.Printf("获取参数失败: %v", err)
		c.JSON(200, Results.ErrBind.Fail("获取参数失败"))
		return
	}

	// 验证验证码是否正确
	if mail.Code == "" {
		c.JSON(200, Results.Err.Fail("验证码不能为空"))
		return
	} else {
		// 从redis中获取验证码
		mailCode, _ := database.GetRedis().Get(context.Background(), mail.Email).Result()
		if mailCode != mail.Code {
			c.JSON(200, Results.Err.Fail("验证码错误"))
			return
		}
	}

	// 验证邮箱是否存在
	userData, _ := models.GetUserByEmail(mail.Email)

	// 邮箱不存在则注册
	if userData.ID == 0 {
		user := models.User{Email: mail.Email, Nickname: mail.Email, Avatar: avatar[rand.Intn(len(avatar))]}
		userId, err := user.AddUser()
		if err != nil {
			c.JSON(200, Results.Err.Fail("注册失败"))
		}

		// 创建token
		token, err := util.GenerateToken(userData.ID, "", mail.Email)
		if err != nil {
			c.JSON(200, Results.ErrToken.Fail("token生成失败"))
			return
		}

		data := make(map[string]interface{})
		data["message"] = "注册并登录成功"
		data["token"] = token
		data["id"] = userId
		data["email"] = mail.Email
		data["nickname"] = mail.Email
		data["avatar"] = userData.Avatar

		//删除redis中的验证码
		database.GetRedis().Del(context.Background(), mail.Email)

		c.JSON(200, Results.Ok.Success(data))
		return
	} else {
		// 邮箱存在则登录
		// 创建token
		token, err := util.GenerateToken(userData.ID, "", mail.Email)
		if err != nil {
			c.JSON(200, Results.ErrToken.Fail("token生成失败"))
			return
		}

		data := make(map[string]interface{})
		data["token"] = token
		data["id"] = userData.ID
		data["email"] = userData.Email
		data["nickname"] = userData.Nickname
		data["avatar"] = userData.Avatar

		//删除redis中的验证码
		database.GetRedis().Del(context.Background(), mail.Email)

		c.JSON(200, Results.Ok.Success(data))
		return
	}

}

// GetPhoneCode 获取手机验证码
func GetPhoneCode(c *gin.Context) {
	var phone util.Phone

	err := c.BindJSON(&phone)
	if err != nil {
		fmt.Printf("获取参数失败: %v", err)
		return
	}

	phoneCode, _ := database.GetRedis().Get(context.Background(), phone.To).Result()

	if phoneCode != "" {
		phone.Code = phoneCode
	} else {
		phoneCode = strconv.Itoa(rand.Intn(9000) + 1000)
		phone.Code = phoneCode
	}

	err = phone.SendSMS() // 发送短信
	if err != nil {
		c.JSON(200, Results.Err.Fail("验证码发送失败"))
		return
	}

	// 将验证码存入redis
	err = database.GetRedis().Set(context.Background(), phone.To, phoneCode, time.Minute*15).Err()
	if err != nil {
		fmt.Printf("验证码存入redis失败: %v", err)
		return
	}

	c.JSON(200, Results.Ok.Success(phone))
}

// GetMailCode 获取邮箱验证码
func GetMailCode(c *gin.Context) {
	var mail util.Mail

	err := c.BindJSON(&mail)
	if err != nil {
		fmt.Printf("获取参数失败: %v", err)
		return
	}

	mail.Subject = "注册验证码"

	mailCode, _ := database.GetRedis().Get(context.Background(), mail.To).Result()
	if mailCode != "" {
		mail.Body = fmt.Sprintf("您的验证码为: %v, 有效期为15分钟。", mailCode)
	} else {
		mailCode = strconv.Itoa(rand.Intn(9000) + 1000)
		mail.Body = fmt.Sprintf("您的验证码为: %v, 有效期为15分钟。", mailCode)
	}

	err = mail.SendMail() // 发送邮件
	if err != nil {
		c.JSON(200, Results.Err.Fail("发送失败"))
		return
	}

	database.GetRedis().Set(context.Background(), mail.To, mailCode, time.Minute*15)

	c.JSON(200, Results.Ok.Success("发送成功"))
}

// UpdateUser 用户更新
func UpdateUser(c *gin.Context) {
	var user map[string]interface{}

	id, _ := c.Get("id")

	err := c.BindJSON(&user)
	if err != nil {
		fmt.Printf("获取参数失败: %v", err)
		c.JSON(200, Results.Err.Fail("获取参数失败"))
		return
	}

	user["update_time"] = time.Now().UnixMilli()

	fmt.Printf("获取参数成功: %v %v", user, id)

	err = models.UpdateUser(id.(int64), user)

	if err != nil {
		c.JSON(200, Results.Err.Fail("更新失败"))
		return
	}

	c.JSON(200, Results.Ok.Success("更新成功"))
}

// UpdateNickname 更新昵称
func UpdateNickname(c *gin.Context) {
	var nickname string

	id, _ := c.Get("id")

	err := c.BindJSON(&nickname)
	if err != nil {
		fmt.Printf("获取参数失败: %v", err)
		c.JSON(200, Results.Err.Fail("获取参数失败"))
		return
	}

	fmt.Printf("获取参数成功: %v %v", nickname, id)

	//判断昵称是否存在
	user, err := models.GetUserByNickname(nickname)

	if user.ID != 0 {
		c.JSON(200, Results.Err.Fail("昵称已存在"))
		return
	}

	err = models.UpdateNickname(id.(int64), nickname)

	if err != nil {
		c.JSON(200, Results.Err.Fail("更新失败"))
		return
	}

	c.JSON(200, Results.Ok.Success("更新成功"))

}
