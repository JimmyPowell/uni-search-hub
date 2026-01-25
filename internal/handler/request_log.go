package handler

import (
	"errors"
	"strconv"
	"time"
	"uni-search-hub/internal/model"
	"uni-search-hub/internal/server"
	"uni-search-hub/pkg/common"

	"github.com/gin-gonic/gin"
)

// GetRequestLogs 查询请求日志（审计）：
// - 普通用户：只能看自己的请求日志
// - Admin/Root：可看全量；可按 user_id 过滤
// 支持分页与筛选（provider / endpoint / status_code / token_id / channel_id / request_id / 时间区间）。
func GetRequestLogs(c *gin.Context) {
	pageInfo := model.GetPageQuery(c)

	role := c.GetInt("role")
	currentUserID := c.GetInt("id")

	// 过滤条件（query params）
	f := model.RequestLogFilter{
		Provider:   c.Query("provider"),
		Endpoint:   c.Query("endpoint"),
		RequestID:  c.Query("request_id"),
		StatusCode: parseInt(c.Query("status_code")),
		TokenID:    parseInt(c.Query("token_id")),
		ChannelID:  parseInt(c.Query("channel_id")),
	}

	// user_id：仅 Admin 及以上允许指定；否则强制为当前用户
	if role >= common.RoleAdminUser {
		f.UserID = parseInt(c.Query("user_id"))
	} else {
		f.UserID = currentUserID
	}

	// 时间区间（毫秒时间戳，精确到秒/毫秒都可）
	startMs, err := parseInt64(c.Query("start_time"))
	if err == nil && startMs > 0 {
		t := time.UnixMilli(startMs)
		f.StartTime = &t
	}
	endMs, err := parseInt64(c.Query("end_time"))
	if err == nil && endMs > 0 {
		t := time.UnixMilli(endMs)
		f.EndTime = &t
	}

	logs, total, err := model.GetRequestLogs(pageInfo, &f)
	if err != nil {
		server.ApiError(c, err)
		return
	}
	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(logs)
	server.ApiSuccess(c, pageInfo)
}

func parseInt(s string) int {
	if s == "" {
		return 0
	}
	v, _ := strconv.Atoi(s)
	return v
}

func parseInt64(s string) (int64, error) {
	if s == "" {
		return 0, errors.New("empty")
	}
	return strconv.ParseInt(s, 10, 64)
}
