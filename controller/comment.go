package controller

import (
	"GinProject/response"
	"GinProject/service"
	"github.com/gin-gonic/gin"
)

type CommentController struct {
}

func Comment() CommentController {
	return CommentController{}
}

func (com CommentController) GetCommentById(c *gin.Context) {
	id := c.Param("id")
	comment := service.GetCommentById(id)
	if comment != nil {
		response.RspSuccess(c, comment)
	} else {
		response.RspError(c, response.CodeCommentNotExist)
	}
}

func (com CommentController) GetCommentByBlog(c *gin.Context) {
	id := c.Param("blog_id")
	comments := service.GetCommentByBlog(id)
	if comments != nil {
		response.RspSuccess(c, comments)
	} else {
		response.RspError(c, response.CodeCommentNotExist)
	}
}
