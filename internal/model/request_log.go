package model

import (
	"time"
	"uni-search-hub/pkg/database"

	"gorm.io/gorm"
)

// RequestLog 审计/计费记录（MVP：仅记录最关键字段）。
type RequestLog struct {
	Id        int            `json:"id" gorm:"primaryKey"`
	RequestId string         `json:"request_id" gorm:"type:varchar(64);index"`

	UserId    int    `json:"user_id" gorm:"index"`
	TokenId   int    `json:"token_id" gorm:"index"`
	ChannelId int    `json:"channel_id" gorm:"index"`
	Provider  string `json:"provider" gorm:"type:varchar(32);index"`

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

