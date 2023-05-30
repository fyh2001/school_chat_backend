package controllers

import (
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"schoolChat/app/models"
	Results "schoolChat/app/result"
	"schoolChat/util"
	"time"
)

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

// Login 用户登录
func Login(c *gin.Context) {
	var user models.UserLogin
	err := c.BindJSON(&user)

	if err != nil {
		fmt.Printf("获取参数失败: %v", err)
	}

	valid := validation.Validation{}

	a := models.UserLogin{Username: user.Username, Password: user.Password}
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := Results.ErrBind.Code //无效参数

	if ok {
		// 去数据库中查询用户是否存在
		userData := models.CheckUser(user.Username)
		if userData.ID > 0 {
			// 验证密码是否正确
			if util.MD5(user.Password) != userData.Password {
				code = Results.ErrPassword.Code //密码错误
			} else {
				// 创建token
				token, err := util.GenerateToken(userData.ID, user.Username, user.Password)
				if err != nil {
					code = Results.ErrToken.Code //token生成失败
				} else {
					data["token"] = token
					data["id"] = userData.ID
					data["username"] = userData.Username
					data["nickname"] = userData.Nickname
					data["avatar"] = userData.Avatar

					code = Results.Ok.Code //成功
				}
			}
		} else {
			code = Results.ErrUser.Code //用户不存在
		}
	} else {
		for _, err := range valid.Errors {
			data[err.Key] = err.Message //将错误信息放入data中
		}
	}

	c.JSON(200, gin.H{
		"code": code,
		"msg":  Results.GetMsg(code),
		"data": data,
	})
}

// Register 用户注册
func Register(c *gin.Context) {
	var user models.User

	err := c.BindJSON(&user)
	if err != nil {
		fmt.Printf("获取参数失败: %v", err)
		return
	}

	userId, err := user.AddUser()
	if err != nil {
		c.JSON(200, Results.Err.Fail("注册失败"))
		return
	}

	c.JSON(200, Results.Ok.Success(userId))
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
