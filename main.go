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

// @description Этот API предоставляет эндпоинты для управления товарами, пользователями и заказами.

// НОВЫЕ СТРОКИ: Описываем схему безопасности Bearer Token (JWT)
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "Для доступа к защищенным эндпоинтам введите 'Bearer ' (с пробелом), а затем ваш JWT. Пример: Bearer eyJhbGciOiJI..."

func main() {
	cfg := config.Load()

	database.InitDB(cfg)
	
	database.CreateInitialAdmin(database.DB, cfg)

	r := router.SetupRouter(cfg)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + cfg.AppPort)
}
