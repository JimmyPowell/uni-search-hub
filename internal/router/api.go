package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetApiRouter(router *gin.Engine) {
	apiRouter := router.Group("/api")
	apiRouter.Use(gzip.Gzip(gzip.DefaultCompression))
	{
		SetSetupRouter(apiRouter)
		SetUserRouter(apiRouter)
		SetTokenRouter(apiRouter)
		SetOptionRouter(apiRouter)
		SetChannelRouter(apiRouter)
		SetSearchRouter(apiRouter)
		SetTavilyRouter(apiRouter)
		SetMetaSoRouter(apiRouter)
		SetZhipuWebSearchRouter(apiRouter)
		SetRequestLogRouter(apiRouter)
		SetMetaRouter(apiRouter)
		SetRedemptionRouter(apiRouter)
	}
}
