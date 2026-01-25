package router

import (
	"uni-search-hub/internal/handler"
	"uni-search-hub/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetMetaRouter(router *gin.RouterGroup) {
	metaRouter := router.Group("/meta")
	metaRouter.Use(middleware.UserAuth())
	{
		metaRouter.GET("/services", handler.GetServiceCatalog)
	}
}
