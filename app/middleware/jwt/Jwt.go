package jwt

import (
	"github.com/gin-gonic/gin"
	"schoolChat/app/result"
	"schoolChat/util"
)

// JWT 自定义中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO
		var code int
		var data interface{}

		code = result.Ok.Code
		token := c.GetHeader("Authorization")
		if token == "" {
			code = result.NoToken.Code // 无token，无权限访问
		} else {
			// 解析token
			claims, err := util.ParseToken(token)
			if err != nil {
				code = result.TokenInvalid.Code // token不合法
			} else {
				id := claims.ID
				email := claims.Email

				// 将当前请求的user信息保存到请求的上下文c上
				c.Set("id", id)
				c.Set("email", email)
			}
			//else if time.Now().Unix() > claims.ExpiresAt {
			//	code = result.TokenExpired.Code // token已过期
			//}

		}
		if code != 200 {
			c.JSON(200, gin.H{
				"code": code,
				"msg":  result.GetMsg(code),
				"data": data,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
