package router

import (
	"uni-search-hub/internal/handler"
	"uni-search-hub/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetSearchRouter(router *gin.RouterGroup) {
	searchRouter := router.Group("/search")
	searchRouter.Use(middleware.TokenAuth())
	{
		searchRouter.POST("/", handler.UnifiedSearch)
	}
}

