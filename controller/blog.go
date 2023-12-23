package controller

import (
	dto "GinProject/dto/blog"
	"GinProject/model"
	"GinProject/response"
	"GinProject/service"
	"github.com/gin-gonic/gin"
)

type BlogController struct {
}

func Blog() BlogController {
	return BlogController{}
}

func (b BlogController) GetBlogById(c *gin.Context) {
	blogId := c.Param("id")
	blog := service.GetBlogById(blogId)
	if blog != nil {
		response.RspSuccess(c, blog)
	} else {
		response.RspError(c, response.CodeBlogNotExist)
	}
}

func (b BlogController) UpdateBlog(c *gin.Context) {
	var blog model.Blog
	err := c.BindJSON(&blog)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	if ok := service.UpdateBlog(&blog); ok {
		response.RspSuccess(c, "update blog success !")
	} else {
		response.RspSuccess(c, "update blog error !")
	}
}

func (b BlogController) AddThumb(c *gin.Context) {
	var thumb model.Thumb
	err := c.BindJSON(&thumb)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	if ok := service.PublishBlogThumb(&thumb); ok {
		response.RspSuccess(c, "add thumb success !")
	} else {
		response.RspSuccess(c, "add thumb error")
	}
}

func (b BlogController) GetThumb(c *gin.Context) {
	blogId := c.Param("blog_id")
	count := service.GetBlogThumb(blogId)
	response.RspSuccess(c, gin.H{
		"thumbs": count,
	})
}

func (b BlogController) GetBlogByUserFavor(c *gin.Context) {
	response.RspSuccess(c, service.GetBlogByUserFavorite(c.Param("user_id")))
}

func (b BlogController) GetBlogByUserId(c *gin.Context) {
	response.RspSuccess(c, service.GetBlogByUserId(c.Param("user_id")))
}

func (b BlogController) PublishBlog(c *gin.Context) {
	blog := &dto.BlogForm{}
	err := c.BindJSON(blog)
	if err != nil {
		response.RspError(c, response.CodeInvalidJson)
		return
	}
	if ok := service.PublishBlog(blog); ok {
		response.RspSuccess(c, "publish blog success !")
		return
	}
	response.RspError(c, response.CodeBlogPublishFail)
}
