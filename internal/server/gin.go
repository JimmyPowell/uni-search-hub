package server

import (
	"net/http"
	"uni-search-hub/internal/config"
	"uni-search-hub/pkg/common"

	"github.com/gin-gonic/gin"
)

func GetGinEngine(config config.ServerConfig) *gin.Engine {
	gin.SetMode(config.Mode)
	r := gin.New()

	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return r
}

func SetContextKey(c *gin.Context, key common.ContextKey, value any) {
	c.Set(string(key), value)
}

func ApiError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"message": err.Error(),
	})
}

func ApiErrorMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"message": msg,
	})
}

func ApiSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    data,
	})
}
