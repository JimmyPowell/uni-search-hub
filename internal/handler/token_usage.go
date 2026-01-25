package handler

import (
	"net/http"
	"strconv"
	"uni-search-hub/internal/model"
	"uni-search-hub/internal/server"
	"uni-search-hub/pkg/common"

	"github.com/gin-gonic/gin"
)

// GetTokenUsageById Session 版 Token 用量查询（管理台使用）：
// - 普通用户只能查看自己的 token
// - Admin/Root 可查看任意 token
func GetTokenUsageById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的参数"})
		return
	}

	role := c.GetInt("role")
	currentUserID := c.GetInt("id")

	var token *model.Token
	var err error

	// 普通用户：直接用 (id, userId) 查询，天然防越权
	if role < common.RoleAdminUser {
		token, err = model.GetTokenByIds(id, currentUserID)
	} else {
		token, err = model.GetTokenById(id)
	}
	if err != nil {
		server.ApiError(c, err)
		return
	}

	expiredAt := token.ExpiredTime
	if expiredAt == -1 {
		expiredAt = 0
	}

	server.ApiSuccess(c, gin.H{
		"object":               "token_usage",
		"name":                 token.Name,
		"total_granted":        token.RemainQuota + token.UsedQuota,
		"total_used":           token.UsedQuota,
		"total_available":      token.RemainQuota,
		"unlimited_quota":      token.UnlimitedQuota,
		"model_limits":         token.GetModelLimitsMap(),
		"model_limits_enabled": token.ModelLimitsEnabled,
		"expires_at":           expiredAt,
	})
}
