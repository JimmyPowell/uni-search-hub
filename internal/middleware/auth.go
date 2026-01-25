package middleware

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"uni-search-hub/internal/model"
	"uni-search-hub/pkg/common"
	"uni-search-hub/pkg/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func validUserInfo(username string, role int) bool {
	// check username is empty
	if strings.TrimSpace(username) == "" {
		return false
	}
	if !common.IsValidateRole(role) {
		return false
	}
	return true
}

func authHelper(c *gin.Context, minRole int) {
	session := sessions.Default(c)
	username := session.Get("username")
	role := session.Get("role")
	id := session.Get("id")
	status := session.Get("status")
	useAccessToken := false
	if username == nil {
		// Check access token
		accessToken := c.Request.Header.Get("Authorization")
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "无权进行此操作，未登录且未提供 access token",
			})
			c.Abort()
			return
		}
		user := model.ValidateAccessToken(accessToken)
		if user != nil && user.Username != "" {
			if !validUserInfo(user.Username, user.Role) {
				c.JSON(http.StatusOK, gin.H{
					"success": false,
					"message": "无权进行此操作，用户信息无效",
				})
				c.Abort()
				return
			}
			// Token is valid
			username = user.Username
			role = user.Role
			id = user.Id
			status = user.Status
			useAccessToken = true
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无权进行此操作，access token 无效",
			})
			c.Abort()
			return
		}
	}
	// get header Uni-Search-Hub-User
	apiUserIdStr := c.Request.Header.Get("Uni-Search-Hub-User")
	if apiUserIdStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "无权进行此操作，未提供 Uni-Search-Hub-User",
		})
		c.Abort()
		return
	}
	apiUserId, err := strconv.Atoi(apiUserIdStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "无权进行此操作，Uni-Search-Hub-User 格式错误",
		})
		c.Abort()
		return

	}
	if id != apiUserId {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "无权进行此操作，Uni-Search-Hub-User 与登录用户不匹配",
		})
		c.Abort()
		return
	}
	if status.(int) == common.UserStatusDisabled {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户已被封禁",
		})
		c.Abort()
		return
	}
	if role.(int) < minRole {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无权进行此操作，权限不足",
		})
		c.Abort()
		return
	}
	if !validUserInfo(username.(string), role.(int)) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无权进行此操作，用户信息无效",
		})
		c.Abort()
		return
	}
	c.Set("username", username)
	c.Set("role", role)
	c.Set("id", id)
	c.Set("group", session.Get("group"))
	c.Set("user_group", session.Get("group"))
	c.Set("use_access_token", useAccessToken)

	c.Next()
}

func TryUserAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		id := session.Get("id")
		if id != nil {
			c.Set("id", id)
		}
		c.Next()
	}
}

func UserAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelper(c, common.RoleCommonUser)
	}
}

func AdminAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelper(c, common.RoleAdminUser)
	}
}

func RootAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelper(c, common.RoleRootUser)
	}
}

func TokenAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 目前兼容两种 Key 的位置
		key := c.Request.Header.Get("Authorization")
		if key == "" {
			key = c.Request.Header.Get("X-Subscription-Token")
		}
		if strings.HasPrefix(key, "Bearer ") || strings.HasPrefix(key, "bearer ") {
			key = strings.TrimSpace(key[7:])
		}

		token, err := model.ValidateUserToken(key)
		if token != nil {
			id := c.GetInt("id")
			if id == 0 {
				c.Set("id", token.UserId)
			}
			// 供后续审计/扣费使用（避免重复查库）
			c.Set("token_id", token.Id)
			c.Set("token", token)
		}
		if err != nil {
			abortWithMessage(c, http.StatusUnauthorized, err.Error())
			return
		}

		allowIps := token.GetIpLimits()
		if len(allowIps) > 0 {
			clientIp := c.ClientIP()
			ip := net.ParseIP(clientIp)
			if ip == nil {
				abortWithMessage(c, http.StatusForbidden, "无法解析客户端 IP 地址")
				return
			}
			if utils.IsIpInCIDRList(ip, allowIps) == false {
				abortWithMessage(c, http.StatusForbidden, "您的 IP 不在令牌允许访问的列表中")
				return
			}
		}

		userCache, err := model.GetUserCache(token.UserId)
		if err != nil {
			abortWithMessage(c, http.StatusInternalServerError, err.Error())
			return
		}
		userEnabled := userCache.Status == common.UserStatusEnabled
		if !userEnabled {
			abortWithMessage(c, http.StatusForbidden, "用户已被封禁")
			return
		}

		userCache.WriteContext(c)

		// TODO 组认证

		c.Next()
	}
}
