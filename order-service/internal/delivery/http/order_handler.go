package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/ecommerce/order-service/internal/entity"
	"github.com/your-username/ecommerce/order-service/internal/usecase"
)

// OrderHandler обрабатывает HTTP-запросы для заказов
type OrderHandler struct {
	orderUseCase *usecase.OrderUseCase
}

// NewOrderHandler создает новый экземпляр OrderHandler
func NewOrderHandler(orderUseCase *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
	}
}

// CreateOrder обрабатывает запрос на создание заказа
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order entity.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	if err := h.orderUseCase.CreateOrder(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании заказа"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrder обрабатывает запрос на получение заказа
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.orderUseCase.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заказ не найден"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateOrder обрабатывает запрос на обновление заказа
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var order entity.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	order.ID = id
	if err := h.orderUseCase.UpdateOrder(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении заказа"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// DeleteOrder обрабатывает запрос на удаление заказа
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	if err := h.orderUseCase.DeleteOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении заказа"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAllOrders обрабатывает запрос на получение всех заказов
func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	orders, err := h.orderUseCase.GetAllOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении списка заказов"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrdersByUserID обрабатывает запрос на получение заказов пользователя
func (h *OrderHandler) GetOrdersByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	orders, err := h.orderUseCase.GetOrdersByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении списка заказов"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
