package router

import (
	"uni-search-hub/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetUserRouter(router *gin.Engine) {
	userRouter := router.Group("/api/user")
	{
		userRouter.GET("/", handler.GetAllUsers)
		userRouter.GET("/search", handler.SearchUsers)
		userRouter.GET("/:id", handler.GetUser)
		userRouter.POST("/", handler.CreateUser)
		userRouter.PUT("/", handler.UpdateUser)
		userRouter.DELETE("/:id", handler.DeleteUser)
	}
}
