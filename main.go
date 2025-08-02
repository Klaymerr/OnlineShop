package main

import (
	"OnlineShop/internal/database"
	"OnlineShop/internal/router"
)

func main() {
	database.InitDB()

	r := router.SetupRouter()

	r.Run(":8080")
}
