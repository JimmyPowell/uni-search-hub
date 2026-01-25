package router

import (
	"uni-search-hub/internal/handler"
	"uni-search-hub/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetRequestLogRouter 请求日志（审计）查询接口：
// - 普通用户：只能看自己的请求日志
// - Admin/Root：可看全量并支持按 user_id 筛选
func SetRequestLogRouter(router *gin.RouterGroup) {
	logRouter := router.Group("/request_log")
	logRouter.Use(middleware.UserAuth())
	{
		logRouter.GET("/", handler.GetRequestLogs)
	}
}
