package model

import (
	"errors"
	"time"
	"uni-search-hub/pkg/database"

	"gorm.io/gorm"
)

var ErrTwoFANotEnabled = errors.New("用户未启用2FA")

// TwoFA 用户2FA设置表
type TwoFA struct {
	Id             int            `json:"id" gorm:"primaryKey"`
	UserId         int            `json:"user_id" gorm:"unique;not null;index"`
	Secret         string         `json:"-" gorm:"type:varchar(255);not null"` // TOTP密钥，不返回给前端
	IsEnabled      bool           `json:"is_enabled"`
	FailedAttempts int            `json:"failed_attempts" gorm:"default:0"`
	LockedUntil    *time.Time     `json:"locked_until,omitempty"`
	LastUsedAt     *time.Time     `json:"last_used_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// TwoFABackupCode 备用码使用记录表
type TwoFABackupCode struct {
	Id        int            `json:"id" gorm:"primaryKey"`
	UserId    int            `json:"user_id" gorm:"not null;index"`
	CodeHash  string         `json:"-" gorm:"type:varchar(255);not null"` // 备用码哈希
	IsUsed    bool           `json:"is_used"`
	UsedAt    *time.Time     `json:"used_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// GetTwoFAByUserId 根据用户ID获取2FA设置
func GetTwoFAByUserId(userId int) (*TwoFA, error) {
	if userId == 0 {
		return nil, errors.New("用户ID不能为空")
	}

	var twoFA TwoFA
	err := database.DB.Where("user_id = ?", userId).First(&twoFA).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 返回nil表示未设置2FA
		}
		return nil, err
	}

	return &twoFA, nil
}

// IsTwoFAEnabled 检查用户是否启用了2FA
func IsTwoFAEnabled(userId int) bool {
	twoFA, err := GetTwoFAByUserId(userId)
	if err != nil || twoFA == nil {
		return false
	}
	return twoFA.IsEnabled
}
