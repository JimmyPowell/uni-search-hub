package router

import (
	"uni-search-hub/internal/handler"
	"uni-search-hub/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetRedemptionRouter(router *gin.RouterGroup) {
	redemptionRouter := router.Group("/redemption")
	redemptionRouter.Use(middleware.AdminAuth())
	{
		redemptionRouter.GET("/", handler.GetAllRedemptions)
		redemptionRouter.GET("/search", handler.SearchRedemptions)
		redemptionRouter.GET("/:id", handler.GetRedemption)
		redemptionRouter.POST("/", handler.AddRedemption)
		redemptionRouter.PUT("/", handler.UpdateRedemption)
		redemptionRouter.DELETE("/invalid", handler.DeleteInvalidRedemption)
		redemptionRouter.DELETE("/:id", handler.DeleteRedemption)
	}
}
