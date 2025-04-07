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

	"github.com/your-username/ecommerce/inventory-service/internal/adapter/mongodb"
	"github.com/your-username/ecommerce/inventory-service/internal/delivery/http"
	"github.com/your-username/ecommerce/inventory-service/internal/usecase"
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
	productRepo := mongodb.NewProductRepository(db)

	// Инициализация usecase
	productUseCase := usecase.NewProductUseCase(productRepo)

	// Инициализация обработчиков
	productHandler := http.NewProductHandler(productUseCase)

	// Настройка роутера
	r := gin.Default()

	// Маршруты для товаров
	products := r.Group("/api/products")
	{
		products.POST("", productHandler.CreateProduct)
		products.GET("", productHandler.GetAllProducts)
		products.GET("/:id", productHandler.GetProduct)
		products.PUT("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
	}

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
