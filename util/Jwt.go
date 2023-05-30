package util

import (
	"github.com/golang-jwt/jwt"
	"time"
)

var jwtSecret = []byte("schoolChat") // jwt密钥

// Claims Claims是一些用户信息状态和额外的jwt参数
type Claims struct {
	ID       int64     `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Time     time.Time `json:"time"`
	jwt.StandardClaims
}

// GenerateToken 根据用户的用户名和密码参数token
func GenerateToken(id int64, username, password string) (string, error) {
	nowTime := time.Now()
	//expireTime := nowTime.Add(3 * time.Hour).Unix() // 过期时间为3小时

	claims := Claims{
		ID:       id,
		Username: username,
		Password: password,
		Time:     nowTime,
		StandardClaims: jwt.StandardClaims{
			//ExpiresAt: expireTime,   // 过期时间
			Issuer: "schoolChat", // 签发人
		},
	}
	// 该方法内部生成签名字符串，再用于获取完整、已签名的token
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 生成完整的token字符串
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken 解析token,根据传入的token值获取到Claims对象信息(进而获取其中的用户名和密码)
func ParseToken(token string) (*Claims, error) {
	// 用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil // 返回签名密钥
	})
	if tokenClaims != nil {
		// 校验token是否有效
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目结构体都是用指针传递，节省空间
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid { // Valid()验证基于时间的声明
			return claims, nil // 返回Claims对象
		}
	}
	return nil, err
}
