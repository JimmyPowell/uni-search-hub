package model

import (
	"uni-search-hub/pkg/database"
	"uni-search-hub/pkg/utils"
)

func MigrateDB() error {
	err := database.DB.AutoMigrate(
		&Channel{},
		&Token{},
		&User{},
		&Option{},
		&RequestLog{},
		&TwoFA{},
		&TwoFABackupCode{},
		&Redemption{},
	)
	if err != nil {
		return err
	}
	// Backfill legacy rows: old code path inserted users without CreatedAt/UpdatedAt fields,
	// leaving DB timestamps NULL (which shows up as 0001-01-01... in JSON).
	if err := database.DB.Exec(
		"UPDATE users SET created_at = NOW(3) WHERE created_at IS NULL OR created_at = '0000-00-00 00:00:00'",
	).Error; err != nil {
		utils.SysLog("failed to backfill users.created_at: " + err.Error())
	}
	if err := database.DB.Exec(
		"UPDATE users SET updated_at = NOW(3) WHERE updated_at IS NULL OR updated_at = '0000-00-00 00:00:00'",
	).Error; err != nil {
		utils.SysLog("failed to backfill users.updated_at: " + err.Error())
	}
	// Backfill request_logs.action for legacy rows.
	if err := database.DB.Exec(
		"UPDATE request_logs SET action = ? WHERE action IS NULL OR action = ''",
		RequestActionProxyRequest,
	).Error; err != nil {
		utils.SysLog("failed to backfill request_logs.action: " + err.Error())
	}
	return nil
}
