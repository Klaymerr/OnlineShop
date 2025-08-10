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

		protectedRoutes.POST("orders", createOrder)
		protectedRoutes.GET("orders", getOrders)
	}

	adminRoutes := r.Group("/")
	adminRoutes.Use(AuthMiddleware(), AdminMiddleware())
	{
		adminRoutes.POST("users/:id/promote", promoteUserToAdmin)

		adminRoutes.POST("products", createProduct)
		adminRoutes.PUT("products/:id", updateProduct)
		adminRoutes.DELETE("products/:id", deleteProduct)
	}

	return r
}
