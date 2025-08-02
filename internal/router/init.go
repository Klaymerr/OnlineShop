package router

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	productRoutes := r.Group("/products")
	{
		productRoutes.GET("", getProducts)
		productRoutes.GET("/:id", getProduct)
		productRoutes.POST("", createProduct)
		productRoutes.PUT("/:id", updateProduct)
		productRoutes.DELETE("/:id", deleteProduct)
	}

	return r
}
