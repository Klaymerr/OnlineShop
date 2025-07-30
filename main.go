package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func getProducts(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"message": "Hello, World!",
		},
	)
}

func getProduct(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"message": "Hello, World!",
		},
	)
}

func createProduct(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"message": "Hello, World!",
		},
	)
}

func updateProduct(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"message": "Hello, World!",
		},
	)
}

func deleteProduct(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"message": "Hello, World!",
		},
	)
}

func main() {
	r := gin.Default()

	r.GET("/products", getProducts)

	r.GET("/products/:id", getProduct)
	r.POST("/products", createProduct)
	r.PUT("/products/:id", updateProduct)
	r.DELETE("/products/:id", deleteProduct)

	r.Run("localhost:8080")
}
