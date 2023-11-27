package controller

import (
	"GinProject/dto/user"
	"GinProject/response"
	"GinProject/service"
	"GinProject/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
}

func User() UserController {
	return UserController{}
}

func (u UserController) Login(c *gin.Context) {
	loginform := &dto.LoginForm{}
	err := c.BindJSON(loginform)
	if err != nil {
		response.RspError(c, response.CodeInvalidParams)
		return
	}
	ok := service.Login(loginform.LoginKey, loginform.Password)
	if ok {
		token, err := utils.GenerateToken(loginform.LoginKey, loginform.Password)
		if err != nil {
			response.RspError(c, response.CodeTokenGenerateFail)
		} else {
			response.RspSuccess(c, gin.H{
				"status": "login success",
				"token":  token,
			})
		}
	} else {
		response.RspError(c, response.CodeInvalidPassword)
	}
}

func (u UserController) SendEmail(c *gin.Context) {
	email := c.PostForm("email")
	if service.SendEmail(email) {
		response.RspSuccess(c, "email sent success !")
	} else {
		response.RspError(c, response.CodeEmailSendFail)
	}
}

func (u UserController) Register(c *gin.Context) {
	form := &dto.RegisterForm{}
	err := c.BindJSON(form)
	if err != nil {
		response.RspError(c, response.CodeInvalidParams)
	}
	if service.Register(form) {
		response.RspSuccess(c, "register success")
	} else {
		response.RspError(c, response.CodeRegisterFail)
	}
}

func (u UserController) Logout(c *gin.Context) {
	userKey := c.PostForm("userKey")
	if ok := service.Logout(userKey); ok {
		response.RspSuccess(c, "user logout success !")
	} else {
		response.RspError(c, response.CodeNotLogin)
	}
}
