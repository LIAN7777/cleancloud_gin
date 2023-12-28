package controller

import (
	"GinProject/model"
	"GinProject/response"
	"GinProject/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type UserMessageController struct {
}

func UserMessage() UserMessageController {
	return UserMessageController{}
}

func (u UserMessageController) GetMessageByUser(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	messages := service.GetMessageByUser(int64(userId))
	if messages != nil {
		response.RspSuccess(c, messages)
	} else {
		response.RspError(c, response.CodeGetMessageFail)
	}
}

func (u UserMessageController) AddUserMessage(c *gin.Context) {
	message := &model.Usermessage{}
	err := c.BindJSON(message)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	res := service.AddUserMessage(message)
	response.RspSuccess(c, res)
}

func (u UserMessageController) DeleteUserMessage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res := service.DeleteUserMessage(int64(id))
	response.RspSuccess(c, res)
}
