package main

import (
	"OnlineShop/config"
	"OnlineShop/internal/database"
	"OnlineShop/internal/router"

	_ "OnlineShop/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API для простого интернет-магазина
// @version 1.0

func main() {
	cfg := config.Load()

	database.InitDB(cfg)

	r := router.SetupRouter(cfg)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + cfg.AppPort)
}
