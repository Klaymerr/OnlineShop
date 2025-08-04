package router

import (
	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	Message string `json:"error" example:"Product not found"`
}

type UpdateProductInput struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"gte=0"`
}

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
