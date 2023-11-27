package controller

import (
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
