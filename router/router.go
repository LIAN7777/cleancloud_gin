package router

import (
	"GinProject/controller"
	"GinProject/middleware"
	"github.com/gin-gonic/gin"
	"time"
)

func Router() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	v1.GET("/", controller.New().Test)
	v1.GET("/city/:city", controller.New().GetCity)
	v1.GET("/user/:userId", controller.New().GetUserById)
	v1.GET("/sign/:id", controller.User().SignIn)
	v1.POST("/get_sign", controller.User().GetUserSign)
	v1.POST("/login", controller.User().Login)
	v1.POST("/logout", controller.User().Logout)
	v1.POST("/email", controller.User().SendEmail)
	v1.POST("/register", controller.User().Register)
	v2 := v1.Group("/auth")
	v2.Use(middleware.AuthMiddleware())
	v2.GET("/test", controller.New().JwtAuthTest)

	blog := v1.Group("/blog")
	blog.GET("/:id", controller.Blog().GetBlogById)
	blog.POST("/update", controller.Blog().UpdateBlog)
	blog.GET("/thumb/:blog_id", controller.Blog().GetThumb)
	blog.POST("/addthumb", controller.Blog().AddThumb)
	blog.GET("/favorite/:user_id", controller.Blog().GetBlogByUserFavor)
	blog.GET("/user/:user_id", controller.Blog().GetBlogByUserId)
	blog.POST("/publish", controller.Blog().PublishBlog)
	blog.GET("/addhits/:id", controller.Blog().AddBlogHits)
	blog.GET("/hot/:limit", controller.Blog().GetHotBlogs)

	admin := v1.Group("/admin")
	admin.POST("/login", controller.Admin().Login)
	admin.GET("/:id", controller.Admin().GetAdminById)
	admin.POST("/update_info", controller.Admin().UpdateAdmin)
	admin.POST("/update_psw", controller.Admin().UpdateAdminPsw)

	comment := v1.Group("/comment")
	comment.GET("/:id", controller.Comment().GetCommentById)
	comment.GET("/blog/:blog_id", controller.Comment().GetCommentByBlog)
	comment.GET("/reported", controller.Comment().GetReportedComment)
	comment.DELETE("/delete/:id", controller.Comment().DeleteCommentById)
	comment.GET("/change_status/:id", controller.Comment().ChangeStatus)
	comment.POST("/publish", controller.Comment().PublishComment)

	follow := v1.Group("/follow")
	follow.GET("/user/:id", controller.Follow().GetFollowByUser)

	//测试用接口
	//限流测试
	limiter := middleware.LimiterMiddleWare(5, 5, time.Second, 2, time.Second)
	limiter2 := middleware.LimiterMiddleWare(2, 2, time.Second, 3, time.Second)
	test := r.Group("test")
	test2 := r.Group("test2")
	test2.Use(limiter2)
	test.Use(limiter)
	test.POST("file", controller.New().FileTranTest)
	test.POST("comment", controller.New().QueueTest)
	test.GET("limiter", controller.New().LimiterTest)
	test2.GET("limiter2", controller.New().LimiterTest)
	return r
}
