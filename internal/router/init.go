package router

import (
	"OnlineShop/config"
	"github.com/gin-gonic/gin"
)

var jwtKey []byte

type HTTPError struct {
	Message string `json:"error" example:"Product not found"`
}

type SuccessMessage struct {
	Message string `json:"message" example:"Product deleted successfully"`
}

func SetupRouter(cfg *config.Config) *gin.Engine {
	jwtKey = cfg.JWTSecretKey

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

		protectedRoutes.POST("orders", createOrder)
	}

	return r
}
