package router

import (
	"OnlineShop/internal/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UpdateProductInput struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"gte=0"`
}

// @Summary      Получить список всех товаров
// @Description  Возвращает массив всех товаров, доступных в магазине
// @Tags         Товары (Products)
// @Produce      json
// @Success      200  {array}   database.Product
// @Router       /products [get]
func getProducts(c *gin.Context) {
	var products []database.Product
	database.DB.Find(&products)
	c.JSON(http.StatusOK, products)
}

// @Summary      Получить товар по ID
// @Description  Получает информацию о конкретном товаре по его ID
// @Tags         Товары (Products)
// @Produce      json
// @Param        id   path      int  true  "ID Товара"
// @Success      200  {object}  database.Product
// @Failure      404  {object}  router.HTTPError
// @Router       /products/{id} [get]
func getProduct(c *gin.Context) {
	id := c.Param("id")
	var product database.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, HTTPError{Message: "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// @Summary      Создать новый товар
// @Description  Добавляет новый товар в базу данных. ID в теле запроса игнорируется.
// @Tags         Товары (Products)
// @Accept       json
// @Produce      json
// @Param        product  body      database.Product  true  "Данные для создания нового товара"
// @Security     BearerAuth
// @Success      201      {object}  database.Product
// @Failure      400      {object}  router.HTTPError
// @Router       /products [post]
func createProduct(c *gin.Context) {
	var product database.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, HTTPError{Message: err.Error()})
		return
	}
	database.DB.Create(&product)
	c.JSON(http.StatusCreated, product)
}

// @Summary      Обновить существующий товар
// @Description  Полностью обновляет информацию о товаре с указанным ID
// @Tags         Товары (Products)
// @Accept       json
// @Produce      json
// @Param        id       path      int                   true  "ID Товара для обновления"
// @Param        product  body      router.UpdateProductInput  true  "Новые данные для товара"
// @Security     BearerAuth
// @Success      200      {object}  database.Product
// @Failure      400      {object}  router.HTTPError
// @Failure      404      {object}  router.HTTPError
// @Router       /products/{id} [put]
func updateProduct(c *gin.Context) {
	id := c.Param("id")

	var product database.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, HTTPError{Message: "Product not found"})
		return
	}

	var input UpdateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, HTTPError{Message: err.Error()})
		return
	}

	product.Name = input.Name
	product.Price = input.Price

	database.DB.Save(&product)

	c.JSON(http.StatusOK, product)
}

// @Summary      Удалить товар
// @Description  Удаляет товар из базы данных по его ID
// @Tags         Товары (Products)
// @Produce      json
// @Param        id   path      int  true  "ID Товара для удаления"
// @Security     BearerAuth
// @Success      200  {object}  router.SuccessMessage
// @Failure      404  {object}  router.HTTPError
// @Failure      500  {object}  router.HTTPError
// @Router       /products/{id} [delete]
func deleteProduct(c *gin.Context) {
	id := c.Param("id")

	var product database.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, HTTPError{Message: "Product not found"})
		return
	}

	if err := database.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, HTTPError{Message: "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, SuccessMessage{Message: "Product deleted successfully"})
}
