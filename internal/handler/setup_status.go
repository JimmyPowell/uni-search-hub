package handler

import (
	"net/http"
	"uni-search-hub/internal/model"

	"github.com/gin-gonic/gin"
)

// GetSetupStatus 返回当前系统是否需要初始化（是否存在 Root 用户）。
func GetSetupStatus(c *gin.Context) {
	rootExists := model.RootUserExists()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"rootExists":     rootExists,
			"setupRequired":  !rootExists,
			"setupCompleted": rootExists,
		},
	})
}
