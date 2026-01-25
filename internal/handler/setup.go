package handler

import (
	"uni-search-hub/internal/model"
	"uni-search-hub/pkg/common"
	"uni-search-hub/pkg/database"

	"github.com/gin-gonic/gin"
)

type SetupRequest struct {
	Username           string `json:"username"`
	Password           string `json:"password"`
	ConfirmPassword    string `json:"confirmPassword"`
	SelfUseModeEnabled bool   `json:"SelfUseModeEnabled"`
	DemoSiteEnabled    bool   `json:"DemoSiteEnabled"`
}

// SetUpRootUser 创建一个 Root 用户，仅测试阶段使用
func SetUpRootUser(c *gin.Context) {
	// Check if root user already exists
	rootExists := model.RootUserExists()

	var req SetupRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "请求参数有误",
		})
		return
	}

	// If root doesn't exist, validate and create admin account
	if !rootExists {
		// Validate username length: max 12 characters to align with model.User validation
		if len(req.Username) > 12 {
			c.JSON(200, gin.H{
				"success": false,
				"message": "用户名长度不能超过12个字符",
			})
			return
		}
		// Validate password
		if req.Password != req.ConfirmPassword {
			c.JSON(200, gin.H{
				"success": false,
				"message": "两次输入的密码不一致",
			})
			return
		}

		if len(req.Password) < 8 {
			c.JSON(200, gin.H{
				"success": false,
				"message": "密码长度至少为8个字符",
			})
			return
		}

		// Create root user.
		// Use model.User.Insert to ensure default fields (especially aff_code) are filled,
		// otherwise MySQL unique index on aff_code may conflict on empty string.
		rootUser := model.User{
			Username:    req.Username,
			Password:    req.Password, // Insert() will hash it
			Role:        common.RoleRootUser,
			Status:      common.UserStatusEnabled,
			DisplayName: "Root User",
		}
		if err := rootUser.Insert(0); err != nil {
			c.JSON(200, gin.H{
				"success": false,
				"message": "创建管理员账号失败: " + err.Error(),
			})
			return
		}
		// Override default quota for root.
		rootUser.Quota = 100000000
		_ = database.DB.Model(&rootUser).Update("quota", rootUser.Quota).Error
		c.JSON(200, gin.H{
			"success": true,
			"message": "创建管理员账号成功",
		})
		return
	}
	c.JSON(200, gin.H{
		"success": false,
		"message": "管理员帐号已经存在",
	})
}

// SetupRootUser 用于前端初始化流程的 Root 创建接口（POST /api/setup/root）。
// 保留 SetUpRootUser 兼容历史接口（GET /api/user/root）。
func SetupRootUser(c *gin.Context) {
	SetUpRootUser(c)
}
