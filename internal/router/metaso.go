package router

import (
	"uni-search-hub/internal/handler"
	"uni-search-hub/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetMetaSoRouter 秘塔（MetaSo）原生接口透传。
func SetMetaSoRouter(router *gin.RouterGroup) {
	metaRouter := router.Group("/metaso")
	metaRouter.Use(middleware.TokenAuth())
	{
		metaRouter.POST("/search", handler.MetaSoSearchProxy)
	}
}
