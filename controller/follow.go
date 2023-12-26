package controller

import (
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
