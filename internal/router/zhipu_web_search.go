package router

import (
	"github.com/gin-gonic/gin"
	"uni-search-hub/internal/handler"
	"uni-search-hub/internal/middleware"
)

// SetZhipuWebSearchRouter 智谱 Web Search 原生接口透传。
func SetZhipuWebSearchRouter(router *gin.RouterGroup) {
	zhipuRouter := router.Group("/zhipu_web_search")
	zhipuRouter.Use(middleware.TokenAuth())
	{
		zhipuRouter.POST("/search", handler.ZhipuWebSearchProxy)
	}
}
