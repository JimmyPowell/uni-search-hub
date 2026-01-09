package router

import (
	"uni-search-hub/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetTokenRouter(router *gin.RouterGroup) {
	tokenRouter := router.Group("/token")
	{
		tokenRouter.GET("/", handler.GetAllTokens)
		tokenRouter.GET("/search", handler.SearchTokens)
		tokenRouter.GET("/:id", handler.GetToken)
		tokenRouter.POST("/", handler.AddToken)
		tokenRouter.PUT("/", handler.UpdateToken)
		tokenRouter.DELETE("/:id", handler.DeleteToken)
		tokenRouter.POST("/batch", handler.DeleteTokenBatch)
	}

	usageRoute := router.Group("api/usage")
	{
		tokenUsageRoute := usageRoute.Group("/token")
		{
			tokenUsageRoute.GET("/", handler.GetTokenUsage)
		}
	}
}
