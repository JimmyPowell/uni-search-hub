package router

import (
	"uni-search-hub/internal/handler"
	"uni-search-hub/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetOptionRouter(router *gin.RouterGroup) {
	optionRouter := router.Group("/option")
	optionRouter.Use(middleware.RootAuth())
	{
		optionRouter.GET("/", handler.GetOptions)
		optionRouter.PUT("/", handler.UpdateOption)
	}
}
