package model

import (
	"time"
	"uni-search-hub/pkg/database"
)

// RequestLogFilter 请求日志查询过滤条件（零值表示不筛选）。
type RequestLogFilter struct {
	UserID     int
	TokenID    int
	ChannelID  int
	Provider   string
	Action     string
	Endpoint   string
	RequestID  string
	StatusCode int

	StartTime *time.Time
	EndTime   *time.Time
}

// GetRequestLogs 分页查询请求日志。
func GetRequestLogs(pageInfo *PageInfo, f *RequestLogFilter) (logs []*RequestLog, total int64, err error) {
	tx := database.DB.Model(&RequestLog{})

	if f != nil {
		if f.UserID != 0 {
			tx = tx.Where("user_id = ?", f.UserID)
		}
		if f.TokenID != 0 {
			tx = tx.Where("token_id = ?", f.TokenID)
		}
		if f.ChannelID != 0 {
			tx = tx.Where("channel_id = ?", f.ChannelID)
		}
		if f.Provider != "" {
			tx = tx.Where("provider = ?", f.Provider)
		}
		if f.Action != "" {
			tx = tx.Where("action = ?", f.Action)
		}
		if f.Endpoint != "" {
			tx = tx.Where("endpoint = ?", f.Endpoint)
		}
		if f.RequestID != "" {
			tx = tx.Where("request_id = ?", f.RequestID)
		}
		if f.StatusCode != 0 {
			tx = tx.Where("status_code = ?", f.StatusCode)
		}
		if f.StartTime != nil {
			tx = tx.Where("created_at >= ?", *f.StartTime)
		}
		if f.EndTime != nil {
			tx = tx.Where("created_at <= ?", *f.EndTime)
		}
	}

	if err = tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err = tx.Order("id desc").
		Limit(pageInfo.GetPageSize()).
		Offset(pageInfo.GetStartIdx()).
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}
