package main

import (
	"uni-search-hub/internal/config"
	"uni-search-hub/internal/model"
	"uni-search-hub/internal/router"
	"uni-search-hub/internal/server"
	"uni-search-hub/pkg/common"
	"uni-search-hub/pkg/database"
)

func main() {
	cfg := config.LoadConfig()

	database.InitDB(cfg.Database.MySQL)
	database.InitRedis(cfg.Database.Redis)

	r := server.GetGinEngine(cfg.Server)
	router.SetRouter(r)

	// 初始化配置
	model.InitOptionMap()

	// 热更新配置
	go model.SyncOptions(common.SyncFrequency)

	err := r.Run(":" + cfg.Server.Port)
	if err != nil {
		return
	}
}
