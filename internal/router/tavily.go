package router

import (
	"uni-search-hub/internal/handler"
	"uni-search-hub/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetTavilyRouter Tavily 原生接口透传。
func SetTavilyRouter(router *gin.RouterGroup) {
	tavilyRouter := router.Group("/tavily")
	// 兼容部分客户端/SDK（例如 @agentic/tavily）：会把 key 放在 JSON body 的 api_key 字段中。
	tavilyRouter.Use(middleware.TokenAuthAllowBodyAPIKey())
	{
		tavilyRouter.POST("/search", handler.TavilySearchProxy)
	}
}
