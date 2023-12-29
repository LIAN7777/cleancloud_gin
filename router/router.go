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
	v2 := v1.Group("/auth")
	v2.Use(middleware.AuthMiddleware())
	v2.GET("/test", controller.New().JwtAuthTest)

	user := v1.Group("/user")
	user.GET("/:id", controller.User().GetUserById)
	user.DELETE("/:id", controller.User().DeleteUser)
	user.GET("/sign/:id", controller.User().SignIn)
	user.POST("/get_sign", controller.User().GetUserSign)
	user.POST("/login", controller.User().Login)
	user.POST("/logout", controller.User().Logout)
	user.POST("/email", controller.User().SendEmail)
	user.POST("/register", controller.User().Register)
	user.GET("/change_status/:id", controller.User().ChangeUserStatus)
	user.GET("/real_name/:id", controller.User().UserRealName)
	user.GET("/auth/:id", controller.User().UserAdminAuth)

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
	comment.GET("/hot/:blog_id", controller.Comment().GetHotComment)
	comment.POST("/add_thumb", controller.Comment().AddCommentThumb)

	follow := v1.Group("/follow")
	follow.GET("/user/:id", controller.Follow().GetFollowByUser)
	follow.POST("/add", controller.Follow().AddFollow)
	follow.POST("/delete", controller.Follow().DeleteFollow)
	follow.POST("/judge", controller.Follow().JudgeFollow)

	favor := v1.Group("/favor")
	favor.POST("/add", controller.Favor().AddFavor)
	favor.POST("/delete", controller.Favor().DeleteFavor)
	favor.POST("/judge", controller.Favor().JudgeFavor)

	userMessage := v1.Group("user_message")
	userMessage.GET("/user/:id", controller.UserMessage().GetMessageByUser)
	userMessage.POST("/add", controller.UserMessage().AddUserMessage)
	userMessage.DELETE("/delete/:id", controller.UserMessage().DeleteUserMessage)

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
