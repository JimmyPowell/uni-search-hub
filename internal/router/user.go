package router

import (
	"uni-search-hub/internal/handler"
	"uni-search-hub/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("/user")
	{
		userRouter.GET("/root", handler.SetUpRootUser)

		userRouter.POST("/register", middleware.TurnstileCheck(), handler.Register)
		userRouter.POST("/login", middleware.TurnstileCheck(), handler.Login)
		userRouter.GET("/logout", handler.Logout)

		selfRouter := userRouter.Group("/")
		selfRouter.Use(middleware.UserAuth())
		{
			selfRouter.GET("/self", handler.GetSelf)
			selfRouter.PUT("/self", handler.UpdateSelf)
			selfRouter.DELETE("/self", handler.DeleteSelf)
			selfRouter.GET("/token", handler.GenerateAccessToken)
		}

		adminRouter := userRouter.Group("/")
		adminRouter.Use(middleware.AdminAuth())
		{
			adminRouter.GET("/", handler.GetAllUsers)
			adminRouter.GET("/search", handler.SearchUsers)
			adminRouter.GET("/:id", handler.GetUser)
			adminRouter.POST("/", handler.CreateUser)
			adminRouter.PUT("/", handler.UpdateUser)
			adminRouter.DELETE("/:id", handler.DeleteUser)
		}

	}
}
