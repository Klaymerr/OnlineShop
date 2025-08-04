package main

import (
	"OnlineShop/internal/database"
	"OnlineShop/internal/router"

	_ "OnlineShop/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API для простого интернет-магазина
// @version 1.0

func main() {
	database.InitDB()

	r := router.SetupRouter()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
