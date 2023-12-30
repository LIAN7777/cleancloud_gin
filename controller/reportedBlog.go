package controller

import (
	dto "GinProject/dto/report"
	"GinProject/response"
	"GinProject/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ReportedBlogController struct {
}

func ReportedBlog() ReportedBlogController {
	return ReportedBlogController{}
}

func (r ReportedBlogController) AddReportedBlog(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	res, err := service.AddReportedBlog(int64(id))
	if err != nil {
		response.RspSuccess(c, gin.H{
			"res":     res,
			"message": err.Error(),
		})
		return
	}
	response.RspSuccess(c, gin.H{
		"res":     res,
		"message": "add reported blog success !",
	})
}

func (r ReportedBlogController) GetReportedById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	blog := service.GetReportedBlogById(int64(id))
	if blog != nil {
		response.RspSuccess(c, blog)
		return
	}
	response.RspError(c, response.CodeBlogNotExist)
}

func (r ReportedBlogController) DeleteReported(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res := service.DeleteReportedBlog(int64(id))
	response.RspSuccess(c, res)
}

func (r ReportedBlogController) AddAssistComment(c *gin.Context) {
	comment := &dto.AssistantComment{}
	err := c.BindJSON(comment)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	res, err := service.AddAssistantComment(comment)
	if err != nil {
		response.RspSuccess(c, gin.H{
			"res":     false,
			"message": err.Error(),
		})
		return
	}
	response.RspSuccess(c, gin.H{
		"res":     res,
		"message": "add comment success !",
	})
}
