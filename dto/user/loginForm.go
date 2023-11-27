package dto

type LoginForm struct {
	LoginKey string `json:"login_key" ,binding:"required"`
	Password string `json:"password" ,binding:"required"`
}
