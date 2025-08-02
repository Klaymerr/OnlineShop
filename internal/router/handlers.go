package router

import (
	"OnlineShop/internal/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getProducts(c *gin.Context) {
	var products []database.Product
	database.DB.Find(&products)
	c.JSON(http.StatusOK, products)
}

func getProduct(c *gin.Context) {
	id := c.Param("id")
	var product database.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func createProduct(c *gin.Context) {
	var product database.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&product)
	c.JSON(http.StatusCreated, product)
}

func updateProduct(c *gin.Context) {
	id := c.Param("id")
	var product database.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Save(&product)
	c.JSON(http.StatusOK, product)
}

func deleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&database.Product{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
