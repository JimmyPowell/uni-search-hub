package model

import (
	"encoding/json"
	"time"
	"uni-search-hub/pkg/database"

	"gorm.io/gorm"
)

const (
	// RequestActionProxyRequest 普通 API 请求（透传/统一搜索等，会产生 cost）。
	RequestActionProxyRequest = "proxy_request"
	// RequestActionChannelTest 渠道测试（不扣费，Root 操作）。
	RequestActionChannelTest = "channel_test"
	// RequestActionTokenCharge 充值/额度调整（预留）。
	RequestActionTokenCharge = "token_charge"
)

// RequestLog 审计/计费记录（MVP：仅记录最关键字段）。
type RequestLog struct {
	Id        int    `json:"id" gorm:"primaryKey"`
	RequestId string `json:"request_id" gorm:"type:varchar(64);index"`

	UserId    int    `json:"user_id" gorm:"index"`
	TokenId   int    `json:"token_id" gorm:"index"`
	ChannelId int    `json:"channel_id" gorm:"index"`
	Provider  string `json:"provider" gorm:"type:varchar(32);index"`

	// Action 审计类型：proxy_request/channel_test/token_charge...
	Action string `json:"action" gorm:"type:varchar(32);index;default:'proxy_request'"`
	// Meta 额外信息（JSON 字符串）。避免频繁加列；仅保存非敏感、体积可控的字段。
	Meta string `json:"meta,omitempty" gorm:"type:longtext"`

	Endpoint   string `json:"endpoint" gorm:"type:varchar(128)"`
	StatusCode int    `json:"status_code"`
	LatencyMs  int64  `json:"latency_ms"`

	// Cost 计费单位（MVP：固定 1）。
	Cost  int    `json:"cost"`
	Error string `json:"error" gorm:"type:varchar(255)"`

	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (l *RequestLog) Insert() error {
	return database.DB.Create(l).Error
}

func (l *RequestLog) SetMeta(v any) {
	if v == nil {
		return
	}
	b, err := json.Marshal(v)
	if err != nil {
		return
	}
	l.Meta = string(b)
}
