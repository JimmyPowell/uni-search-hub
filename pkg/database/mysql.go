package database

import (
	"log"
	"time"
	"uni-search-hub/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(config config.MySQLConfig) {
	db, err := gorm.Open(mysql.Open(config.DSN))
	if err != nil {
		log.Fatal(err)
	}

	// 获取通用数据库对象 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置了连接可复用的最大时间

	// 测试连接
	if err = sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	DB = db
}
