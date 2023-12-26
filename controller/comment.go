package controller

import (
	dto "GinProject/dto/user"
	"GinProject/response"
	"GinProject/service"
	"github.com/gin-gonic/gin"
	"strconv"
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

func (com CommentController) GetReportedComment(c *gin.Context) {
	response.RspSuccess(c, service.GetReportedComment())
}

func (com CommentController) DeleteCommentById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if ok := service.DeleteCommentById(int64(id)); ok {
		response.RspSuccess(c, "delete success !")
	} else {
		response.RspError(c, response.CodeCommentNotExist)
	}
}

func (com CommentController) ChangeStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if ok := service.ChangeStatus(int64(id)); ok {
		response.RspSuccess(c, "change success !")
	} else {
		response.RspError(c, response.CodeCommentNotExist)
	}
}

func (com CommentController) PublishComment(c *gin.Context) {
	comment := &dto.CommentForm{}
	err := c.BindJSON(comment)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	if ok := service.PublishComment(comment); ok {
		response.RspSuccess(c, "add comment success")
		return
	}
	response.RspError(c, response.CodeServerBusy)
}
