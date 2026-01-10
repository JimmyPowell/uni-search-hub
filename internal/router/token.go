package router

import (
	"uni-search-hub/internal/handler"
	"uni-search-hub/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetTokenRouter(router *gin.RouterGroup) {
	tokenRouter := router.Group("/token")
	tokenRouter.Use(middleware.UserAuth())
	{
		tokenRouter.GET("/", handler.GetAllTokens)
		tokenRouter.GET("/search", handler.SearchTokens)
		tokenRouter.GET("/:id", handler.GetToken)
		tokenRouter.POST("/", handler.AddToken)
		tokenRouter.PUT("/", handler.UpdateToken)
		tokenRouter.DELETE("/:id", handler.DeleteToken)
		tokenRouter.POST("/batch", handler.DeleteTokenBatch)
	}

	usageRoute := router.Group("/usage")
	{
		tokenUsageRoute := usageRoute.Group("/token")
		tokenUsageRoute.Use(middleware.TokenAuth())
		{
			tokenUsageRoute.GET("/", handler.GetTokenUsage)
		}
	}
}
