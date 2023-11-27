package middleware

import (
	"GinProject/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 密钥验证
			return []byte("LIANSecretKey"), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		//读取redis验证token是否过期
		rLoginKey := "login:jwt:" + token.Claims.(jwt.MapClaims)["loginKey"].(string)
		rTokenStr, _ := utils.Client.Get(rLoginKey).Result()
		if rTokenStr == "" || rTokenStr != tokenString {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		} else {
			utils.Client.Expire(rLoginKey, time.Minute*60)
		}

		////token有效则延长token时效
		//exp := time.Now().Add(time.Second * 30).Unix()
		//// 构造新的claims
		//newClaims := jwt.MapClaims{
		//	"loginKey": token.Claims.(jwt.MapClaims)["loginKey"],
		//	"password": token.Claims.(jwt.MapClaims)["password"],
		//	"exp":      exp,
		//}
		//
		//// 生成新的token
		//newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
		//
		//// 签名并获取完整的token字符串
		//key := []byte("LIANSecretKey")
		//newTokenString, err := newToken.SignedString(key)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//// 将解析后的token保存到上下文中，以便后续处理函数使用
		//c.Set("token", newToken)
		//c.Set("newTokenString", newTokenString)
		c.Next()
	}
}
