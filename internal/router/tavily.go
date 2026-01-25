package router

import (
	"uni-search-hub/internal/handler"
	"uni-search-hub/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetTavilyRouter Tavily 原生接口透传。
func SetTavilyRouter(router *gin.RouterGroup) {
	tavilyRouter := router.Group("/tavily")
	tavilyRouter.Use(middleware.TokenAuth())
	{
		tavilyRouter.POST("/search", handler.TavilySearchProxy)
	}
}

