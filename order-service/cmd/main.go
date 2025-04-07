package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/your-username/ecommerce/order-service/internal/adapter/mongodb"
	"github.com/your-username/ecommerce/order-service/internal/delivery/http"
	"github.com/your-username/ecommerce/order-service/internal/usecase"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Подключение к MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(os.Getenv("MONGODB_DATABASE"))

	// Инициализация репозитория
	orderRepo := mongodb.NewOrderRepository(db)

	// Инициализация usecase
	orderUseCase := usecase.NewOrderUseCase(orderRepo)

	// Инициализация обработчиков
	orderHandler := http.NewOrderHandler(orderUseCase)

	// Настройка роутера
	r := gin.Default()

	// Маршруты для заказов
	orders := r.Group("/api/orders")
	{
		orders.POST("", orderHandler.CreateOrder)
		orders.GET("", orderHandler.GetAllOrders)
		orders.GET("/:id", orderHandler.GetOrder)
		orders.PUT("/:id", orderHandler.UpdateOrder)
		orders.DELETE("/:id", orderHandler.DeleteOrder)
		orders.GET("/user/:user_id", orderHandler.GetOrdersByUserID)
	}

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
