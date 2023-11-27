package router

import (
	"GinProject/controller"
	"GinProject/middleware"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	v1.GET("/", controller.New().Test)
	v1.GET("/city/:city", controller.New().GetCity)
	v1.GET("/user/:userId", controller.New().GetUserById)
	v1.POST("/login", controller.User().Login)
	v1.POST("/logout", controller.User().Logout)
	v1.POST("/email", controller.User().SendEmail)
	v1.POST("/register", controller.User().Register)
	v2 := v1.Group("/auth")
	v2.Use(middleware.AuthMiddleware())
	v2.GET("/test", controller.New().JwtAuthTest)

	blog := v1.Group("blog")
	blog.GET("/:id", controller.Blog().GetBlogById)
	blog.POST("/update", controller.Blog().UpdateBlog)

	//测试用接口
	test := r.Group("test")
	test.POST("file", controller.New().FileTranTest)
	return r
}
