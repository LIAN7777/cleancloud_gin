package controller

import (
	"GinProject/model"
	"GinProject/response"
	"GinProject/service"
	"github.com/gin-gonic/gin"
)

type FavorController struct {
}

func Favor() FavorController {
	return FavorController{}
}

func (f FavorController) AddFavor(c *gin.Context) {
	favor := &model.Favorite{}
	err := c.BindJSON(favor)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	res := service.AddFavorite(favor)
	response.RspSuccess(c, res)
}

func (f FavorController) DeleteFavor(c *gin.Context) {
	favor := &model.Favorite{}
	err := c.BindJSON(favor)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	res := service.DeleteFavorite(favor)
	response.RspSuccess(c, res)
}

func (f FavorController) JudgeFavor(c *gin.Context) {
	favor := &model.Favorite{}
	err := c.BindJSON(favor)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	res := service.JudgeFavorite(favor)
	response.RspSuccess(c, res)
}
