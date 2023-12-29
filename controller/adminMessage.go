package controller

import (
	"GinProject/model"
	"GinProject/response"
	"GinProject/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type AdminMessageController struct {
}

func AdminMessage() AdminMessageController {
	return AdminMessageController{}
}

func (a AdminMessageController) GetMessageByAdmin(c *gin.Context) {
	adminId, _ := strconv.Atoi(c.Param("id"))
	messages := service.GetMessageByAdmin(int64(adminId))
	if messages != nil {
		response.RspSuccess(c, messages)
	} else {
		response.RspError(c, response.CodeGetMessageFail)
	}
}

func (a AdminMessageController) AddAdminMessage(c *gin.Context) {
	message := &model.Adminmessage{}
	err := c.BindJSON(message)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	res := service.AddAdminMessage(message)
	response.RspSuccess(c, res)
}
func (a AdminMessageController) DeleteAdminMessage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res := service.DeleteAdminMessage(int64(id))
	response.RspSuccess(c, res)
}
