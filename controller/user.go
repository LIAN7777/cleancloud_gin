package controller

import (
	"GinProject/dto/user"
	"GinProject/response"
	"GinProject/service"
	"GinProject/utils"
	"github.com/gin-gonic/gin"
	"strconv"
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
	ok := service.UserLogin(loginform.LoginKey, loginform.Password)
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
	if ok := service.UserLogout(userKey); ok {
		response.RspSuccess(c, "user logout success !")
	} else {
		response.RspError(c, response.CodeNotLogin)
	}
}

func (u UserController) SignIn(c *gin.Context) {
	id := c.Param("id")
	if ok := service.UserSignIn(id); ok {
		response.RspSuccess(c, "sign in success !")
	} else {
		response.RspError(c, response.CodeInternalError)
	}
}

func (u UserController) GetUserSign(c *gin.Context) {
	id := c.PostForm("id")
	day, _ := strconv.Atoi(c.PostForm("day"))
	res := service.GetUserSign(id, int64(day))
	response.RspSuccess(c, res)
}

func (u UserController) GetUserById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user := service.GetUserById(int64(id))
	if user != nil {
		response.RspSuccess(c, user)
	} else {
		response.RspError(c, response.CodeUserNotExist)
	}
}

func (u UserController) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res := service.DeleteUser(int64(id))
	response.RspSuccess(c, res)
}

func (u UserController) ChangeUserStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res := service.ChangeUserStatus(int64(id))
	response.RspSuccess(c, res)
}

func (u UserController) UserRealName(c *gin.Context) {
	realName := &dto.RealName{}
	err := c.BindJSON(realName)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	res := service.UserRealName(realName)
	response.RspSuccess(c, res)
}

func (u UserController) UserAdminAuth(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res := service.UserAdminAuth(int64(id))
	response.RspSuccess(c, res)
}

func (u UserController) UpdateUserInfo(c *gin.Context) {
	form := &dto.UserUpdateForm{}
	err := c.BindJSON(form)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	res := service.UpdateUserInfo(form)
	response.RspSuccess(c, res)
}

func (u UserController) UpdateUserPsw(c *gin.Context) {
	form := &dto.UserPswForm{}
	err := c.BindJSON(form)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	res := service.UpdateUserPsw(form)
	response.RspSuccess(c, res)
}
