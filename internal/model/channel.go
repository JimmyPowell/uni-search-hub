package model

import (
	"errors"
	"time"
	"uni-search-hub/pkg/database"

	"gorm.io/gorm"
)

// Channel 渠道：绑定一个 provider + 一个 API key，用于对外请求时做轮询/负载。
type Channel struct {
	Id       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"type:varchar(64);not null;index"`
	Provider string `json:"provider" gorm:"type:varchar(32);not null;index"`

	// ApiKey 属于敏感信息：写入数据库用于调用上游，但默认响应会清空。
	ApiKey string `json:"api_key" gorm:"type:varchar(255);not null"`

	Enabled bool `json:"enabled" gorm:"default:true;index"`
	Weight  int  `json:"weight" gorm:"default:1"`

	// 失败计数：MVP 用于简单熔断/观测（暂不参与选路）。
	FailCount  int  `json:"fail_count" gorm:"default:0"`
	LastUsedAt int64 `json:"last_used_at" gorm:"bigint;default:0"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (c *Channel) Clean() {
	c.ApiKey = ""
}

func GetChannelById(id int) (*Channel, error) {
	if id == 0 {
		return nil, errors.New("id 为空！")
	}
	var ch Channel
	err := database.DB.First(&ch, "id = ?", id).Error
	return &ch, err
}

func GetChannels(pageInfo *PageInfo) (channels []*Channel, total int64, err error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, 0, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err = tx.Model(&Channel{}).Count(&total).Error
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}
	err = tx.Order("id desc").Limit(pageInfo.GetPageSize()).Offset(pageInfo.GetStartIdx()).Find(&channels).Error
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}
	if err = tx.Commit().Error; err != nil {
		return nil, 0, err
	}
	return channels, total, nil
}

func GetEnabledChannelsByProvider(provider string) ([]*Channel, error) {
	if provider == "" {
		return nil, errors.New("provider 为空！")
	}
	var channels []*Channel
	err := database.DB.Where("provider = ? AND enabled = ?", provider, true).Order("id asc").Find(&channels).Error
	return channels, err
}

func (c *Channel) Insert() error {
	if c.Name == "" || c.Provider == "" || c.ApiKey == "" {
		return errors.New("name/provider/api_key 不能为空")
	}
	return database.DB.Create(c).Error
}

func (c *Channel) Update() error {
	if c.Id == 0 {
		return errors.New("id 为空！")
	}
	// 允许不更新 api_key（前端留空表示不修改）
	updates := map[string]interface{}{
		"name":       c.Name,
		"provider":   c.Provider,
		"enabled":    c.Enabled,
		"weight":     c.Weight,
		"fail_count": c.FailCount,
		"last_used_at": c.LastUsedAt,
	}
	if c.ApiKey != "" {
		updates["api_key"] = c.ApiKey
	}
	return database.DB.Model(&Channel{}).Where("id = ?", c.Id).Updates(updates).Error
}

func (c *Channel) Delete() error {
	if c.Id == 0 {
		return errors.New("id 为空！")
	}
	return database.DB.Delete(&Channel{}, "id = ?", c.Id).Error
}

