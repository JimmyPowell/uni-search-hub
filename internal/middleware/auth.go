package middleware

import (
	"bytes"
	"fmt"
	"io"
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
	utils.SysLog(fmt.Sprintf("role: %d", role))
	utils.SysLog(fmt.Sprintf("minRole: %d", minRole))
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
	// Uni-Search-Hub-User is optional. When provided, it must match the authenticated user.
	// This keeps proxy/on-behalf-of flows possible without breaking normal browser session usage.
	apiUserIdStr := c.Request.Header.Get("Uni-Search-Hub-User")
	if apiUserIdStr != "" {
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
		key, source := normalizeTokenKey(
			c.Request.Header.Get("Authorization"),
			c.Request.Header.Get("X-Subscription-Token"),
		)

		token, err := model.ValidateUserToken(key)
		if common.DebugEnabled {
			// Do not log full token. Only log minimal masked info for debugging.
			masked := maskTokenKey(key)
			utils.SysLog(fmt.Sprintf("[TokenAuth] path=%s source=%s token=%s err=%v", c.Request.URL.Path, source, masked, err))
		}
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

// TokenAuthAllowBodyAPIKey 兼容部分 SDK：若未提供 header token，则尝试从 JSON body 的 api_key 字段读取。
// 注意：该 token 仅用于本系统鉴权，业务 handler 需要在转发上游前剥离 api_key，避免泄露。
func TokenAuthAllowBodyAPIKey() func(c *gin.Context) {
	return func(c *gin.Context) {
		key, source := normalizeTokenKey(
			c.Request.Header.Get("Authorization"),
			c.Request.Header.Get("X-Subscription-Token"),
		)

		// Header 未提供时，尝试从 body 读取 api_key（仅限 JSON）。
		if key == "" && strings.Contains(strings.ToLower(c.GetHeader("Content-Type")), "application/json") {
			if c.Request.Body != nil {
				raw, err := io.ReadAll(c.Request.Body)
				if err == nil {
					// Restore body for downstream handler.
					c.Request.Body = io.NopCloser(bytes.NewReader(raw))
					if k := extractJSONAPIKey(raw); k != "" {
						key = k
						source = "body-api-key"
					}
				}
			}
		}

		token, err := model.ValidateUserToken(key)
		if common.DebugEnabled {
			masked := maskTokenKey(key)
			utils.SysLog(fmt.Sprintf("[TokenAuth] path=%s source=%s token=%s err=%v", c.Request.URL.Path, source, masked, err))
		}
		if token != nil {
			id := c.GetInt("id")
			if id == 0 {
				c.Set("id", token.UserId)
			}
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
		c.Next()
	}
}

func normalizeTokenKey(authHeader string, xSubToken string) (key string, source string) {
	key = strings.TrimSpace(authHeader)
	source = "authorization"
	if key == "" {
		key = strings.TrimSpace(xSubToken)
		source = "x-subscription-token"
	}
	if strings.HasPrefix(key, "Bearer ") || strings.HasPrefix(key, "bearer ") {
		key = strings.TrimSpace(key[7:])
	}
	// Compatibility: allow "sk-" prefix if it looks like a prefixed 48-char key.
	// Avoid stripping when the actual key itself happens to start with "sk-".
	if strings.HasPrefix(key, "sk-") && len(key) == 51 {
		key = key[3:]
	}
	return key, source
}

func extractJSONAPIKey(raw []byte) string {
	// Lightweight parse to avoid binding structs here; only need api_key string.
	// If parsing fails, return empty.
	// NOTE: do not log raw body here.
	var m map[string]any
	if err := utils.Unmarshal(raw, &m); err != nil {
		return ""
	}
	v, ok := m["api_key"]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	// Reuse the same sk- compatibility as header path.
	if strings.HasPrefix(s, "sk-") && len(s) == 51 {
		s = s[3:]
	}
	return s
}

func maskTokenKey(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return "<empty>"
	}
	if len(key) <= 6 {
		return "***"
	}
	return key[:3] + "***" + key[len(key)-3:]
}
