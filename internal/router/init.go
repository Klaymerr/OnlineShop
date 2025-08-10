package router

import (
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("JWT_KEY")

type HTTPError struct {
	Message string `json:"error" example:"Product not found"`
}

type SuccessMessage struct {
	Message string `json:"message" example:"Product deleted successfully"`
}

type UpdateProductInput struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"gte=0"`
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	publicRoutes := r.Group("/")
	{
		publicRoutes.GET("products", getProducts)
		publicRoutes.GET("products/:id", getProduct)

		publicRoutes.POST("users/login", loginUser)
		publicRoutes.POST("users/register", registerUser)
	}

	protectedRoutes := r.Group("/")
	protectedRoutes.Use(AuthMiddleware())
	{
		protectedRoutes.GET("users/me", SayHello)

		protectedRoutes.POST("products", createProduct)
		protectedRoutes.PUT("products/:id", updateProduct)
		protectedRoutes.DELETE("products/:id", deleteProduct)
	}

	return r
}
