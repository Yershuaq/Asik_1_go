package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Инициализация роутера
	r := gin.Default()

	// Настройка маршрутов
	setupRoutes(r)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func setupRoutes(r *gin.Engine) {
	// Маршруты для товаров
	products := r.Group("/api/products")
	{
		products.POST("", proxyToInventoryService)
		products.GET("", proxyToInventoryService)
		products.GET("/:id", proxyToInventoryService)
		products.PUT("/:id", proxyToInventoryService)
		products.DELETE("/:id", proxyToInventoryService)
	}

	// Маршруты для заказов
	orders := r.Group("/api/orders")
	{
		orders.POST("", proxyToOrderService)
		orders.GET("", proxyToOrderService)
		orders.GET("/:id", proxyToOrderService)
		orders.PUT("/:id", proxyToOrderService)
		orders.DELETE("/:id", proxyToOrderService)
		orders.GET("/user/:user_id", proxyToOrderService)
	}
}

func proxyToInventoryService(c *gin.Context) {
	// TODO: Реализовать проксирование запросов к inventory-service
	c.JSON(501, gin.H{"error": "Not implemented"})
}

func proxyToOrderService(c *gin.Context) {
	// TODO: Реализовать проксирование запросов к order-service
	c.JSON(501, gin.H{"error": "Not implemented"})
}
