package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/ecommerce/inventory-service/internal/entity"
	"github.com/your-username/ecommerce/inventory-service/internal/usecase"
)

// ProductHandler обрабатывает HTTP-запросы для товаров
type ProductHandler struct {
	productUseCase *usecase.ProductUseCase
}

// NewProductHandler создает новый экземпляр ProductHandler
func NewProductHandler(productUseCase *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
	}
}

// CreateProduct обрабатывает запрос на создание товара
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product entity.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	if err := h.productUseCase.CreateProduct(c.Request.Context(), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании товара"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct обрабатывает запрос на получение товара
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	product, err := h.productUseCase.GetProductByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct обрабатывает запрос на обновление товара
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product entity.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	product.ID = id
	if err := h.productUseCase.UpdateProduct(c.Request.Context(), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении товара"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct обрабатывает запрос на удаление товара
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := h.productUseCase.DeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении товара"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAllProducts обрабатывает запрос на получение всех товаров
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.productUseCase.GetAllProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении списка товаров"})
		return
	}

	c.JSON(http.StatusOK, products)
}
