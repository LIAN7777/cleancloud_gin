package controller

import (
	"GinProject/model"
	"GinProject/response"
	"GinProject/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type FollowController struct {
}

func Follow() FollowController {
	return FollowController{}
}

func (f FollowController) GetFollowByUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	follows := service.GetFollowByUser(int64(id))
	response.RspSuccess(c, follows)
}

func (f FollowController) AddFollow(c *gin.Context) {
	follow := &model.Follow{}
	err := c.BindJSON(follow)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	res := service.AddFollow(follow)
	response.RspSuccess(c, res)
}

func (f FollowController) DeleteFollow(c *gin.Context) {
	follow := &model.Follow{}
	err := c.BindJSON(follow)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	res := service.DeleteFollow(follow)
	response.RspSuccess(c, res)
}

func (f FollowController) JudgeFollow(c *gin.Context) {
	follow := &model.Follow{}
	err := c.BindJSON(follow)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	res := service.JudgeFollow(follow)
	response.RspSuccess(c, res)
}
