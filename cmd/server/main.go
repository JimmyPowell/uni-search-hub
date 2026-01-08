package main

import (
	"uni-search-hub/internal/config"
	"uni-search-hub/internal/router"
	"uni-search-hub/internal/server"
	"uni-search-hub/pkg/database"
)

func main() {
	cfg := config.LoadConfig()

	database.InitDB(cfg.Database.MySQL)
	database.InitRedis(cfg.Database.Redis)

	r := server.GetGinEngine(cfg.Server)
	router.SetRouter(r)

	err := r.Run(":" + cfg.Server.Port)
	if err != nil {
		return
	}
}
