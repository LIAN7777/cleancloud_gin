package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GenerateToken(loginKey string, password string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["loginKey"] = loginKey
	claims["password"] = password
	// 设置过期时间为10天
	claims["exp"] = time.Now().Add(time.Hour * 240).Unix()

	tokenString, err := token.SignedString([]byte("LIANSecretKey"))
	if err != nil {
		return "", err
	}
	Client.Set("login:jwt:"+loginKey, tokenString, time.Minute*60)

	return tokenString, nil
}
