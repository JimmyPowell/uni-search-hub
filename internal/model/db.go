package model

import "uni-search-hub/pkg/database"

func MigrateDB() error {
	err := database.DB.AutoMigrate(
		&Token{},
		&User{},
		&Option{},
	)
	if err != nil {
		return err
	}
	return nil
}
