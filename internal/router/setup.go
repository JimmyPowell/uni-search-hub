package router

import (
	"uni-search-hub/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetSetupRouter 系统初始化相关接口（Root 引导）。
func SetSetupRouter(router *gin.RouterGroup) {
	setupRouter := router.Group("/setup")
	{
		setupRouter.GET("/status", handler.GetSetupStatus)
		setupRouter.POST("/root", handler.SetupRootUser)
	}
}
