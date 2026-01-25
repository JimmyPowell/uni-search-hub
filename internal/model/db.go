package model

import "uni-search-hub/pkg/database"

func MigrateDB() error {
	err := database.DB.AutoMigrate(
		&Channel{},
		&Token{},
		&User{},
		&Option{},
		&RequestLog{},
		&TwoFA{},
		&TwoFABackupCode{},
	)
	if err != nil {
		return err
	}
	return nil
}
