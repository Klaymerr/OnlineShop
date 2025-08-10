package router

import (
	"OnlineShop/internal/database"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type CreateOrderItemInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

type CreateOrderInput struct {
	Items []CreateOrderItemInput `json:"items" binding:"required,min=1"`
}

// @Summary      Создать новый заказ
// @Description  Создает новый заказ для аутентифицированного пользователя. Требует список ID товаров и их количество.
// @Tags         Заказы (Orders)
// @Accept       json
// @Produce      json
// @Param        order  body      CreateOrderInput  true  "Данные для создания нового заказа"
// @Security     BearerAuth
// @Success      201  {object}  database.Order "Возвращает созданный заказ со всеми позициями"
// @Failure      400  {object}  HTTPError      "Ошибка валидации входных данных"
// @Failure      401  {object}  HTTPError      "Ошибка аутентификации"
// @Failure      404  {object}  HTTPError      "Один или несколько товаров не найдены"
// @Failure      500  {object}  HTTPError      "Внутренняя ошибка сервера"
// @Router       /orders [post]
func createOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, HTTPError{Message: "user ID not found in context"})
		return
	}

	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, HTTPError{Message: err.Error()})
		return
	}

	orderToCreate := database.Order{
		CustomerID: userID.(uint),
		OrderDate:  time.Now(),
		Status:     "Pending",
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&orderToCreate).Error; err != nil {
			return err
		}

		for _, itemInput := range input.Items {
			var product database.Product
			if err := tx.First(&product, itemInput.ProductID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("product not found")
				}
				return err
			}

			orderItem := database.OrderItem{
				OrderID:   orderToCreate.ID,
				ProductID: itemInput.ProductID,
				Quantity:  itemInput.Quantity,
				Price:     product.Price,
			}

			if err := tx.Create(&orderItem).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, HTTPError{Message: "One or more products not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, HTTPError{Message: "Failed to create order"})
		return
	}

	var finalOrder database.Order
	if err := database.DB.Preload("Items.Product").First(&finalOrder, orderToCreate.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, HTTPError{Message: "Failed to fetch created order"})
		return
	}

	c.JSON(http.StatusCreated, finalOrder)
}

// @Summary      Получить список заказов пользователя
// @Description  Возвращает все заказы, сделанные аутентифицированным пользователем, с полной информацией о товарах.
// @Tags         Заказы (Orders)
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   database.Order "Массив заказов пользователя"
// @Failure      401  {object}  HTTPError      "Ошибка аутентификации"
// @Failure      500  {object}  HTTPError      "Внутренняя ошибка сервера"
// @Router       /orders [get]
func getOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, HTTPError{Message: "user ID not found in context"})
		return
	}

	var orders []database.Order

	err := database.DB.
		Preload("Items.Product").
		Where("customer_id = ?", userID).
		Order("order_date DESC").
		Find(&orders).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, HTTPError{Message: "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// @Summary      Получить список всех незавершенных заказов
// @Description  Возвращает список всех заказов в статусе "Pending". Доступно только для администраторов.
// @Tags         Администрирование (Admin)
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   database.Order "Список незавершенных заказов"
// @Failure      500  {object}  HTTPError      "Внутренняя ошибка сервера"
// @Router       /orders/pending [get]
func getPendingOrders(c *gin.Context) {
	var orders []database.Order

	err := database.DB.
		Preload("Items.Product").
		Preload("Customer").
		Where("status = ?", "Pending").
		Order("order_date ASC").
		Find(&orders).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, HTTPError{Message: "Failed to fetch pending orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
