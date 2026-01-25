package router

import (
	"uni-search-hub/internal/handler"
	"uni-search-hub/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetChannelRouter(router *gin.RouterGroup) {
	channelRouter := router.Group("/channel")
	channelRouter.Use(middleware.RootAuth())
	{
		channelRouter.GET("/", handler.GetChannels)
		channelRouter.POST("/", handler.CreateChannel)
		channelRouter.PUT("/", handler.UpdateChannel)
		channelRouter.POST("/:id/test", handler.TestChannel)
		channelRouter.DELETE("/:id", handler.DeleteChannel)
	}
}
