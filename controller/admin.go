package controller

import (
	dto "GinProject/dto/admin"
	"GinProject/response"
	"GinProject/service"
	"GinProject/utils"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
}

func Admin() AdminController {
	return AdminController{}
}

func (a AdminController) Login(c *gin.Context) {
	loginForm := &dto.AdminLoginForm{}
	err := c.BindJSON(loginForm)
	if err != nil {
		response.RspError(c, response.CodeInvalidParams)
		return
	}
	if ok := service.AdminLogin(loginForm); ok {
		token, err := utils.GenerateToken(loginForm.AdminId, loginForm.Password)
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

func (a AdminController) GetAdminById(c *gin.Context) {
	id := c.Param("id")
	admin := service.GetAdminById(id)
	if admin == nil {
		response.RspError(c, response.CodeAdminNotExist)
		return
	} else {
		response.RspSuccess(c, admin)
	}
}

func (a AdminController) Logout(c *gin.Context) {
	id := c.PostForm("id")
	if ok := service.AdminLogout(id); ok {
		response.RspSuccess(c, "admin logout success !")
	} else {
		response.RspError(c, response.CodeNotLogin)
	}
}

func (a AdminController) UpdateAdmin(c *gin.Context) {
	form := &dto.AdminUpdateForm{}
	err := c.BindJSON(form)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	if ok := service.UpdateAdminInfo(form); ok {
		response.RspSuccess(c, "update success !")
		return
	} else {
		response.RspError(c, response.CodeUpdateError)
	}
}

func (a AdminController) UpdateAdminPsw(c *gin.Context) {
	form := &dto.AdminPswForm{}
	err := c.BindJSON(form)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	if ok := service.UpdateAdminPassword(form); ok {
		response.RspSuccess(c, "update success !")
		return
	}
	response.RspError(c, response.CodeUpdateError)
}
