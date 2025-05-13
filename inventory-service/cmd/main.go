package main

import (
	"context"
	"log"
	"net"
	"os"
	"time" // Добавляем импорт time

	// Другие импорты...
	"github.com/Yershuaq/Asik_1_go/inventory-service/internal/adapter/cache" // Импорт адаптера кэша
	"github.com/Yershuaq/Asik_1_go/inventory-service/internal/adapter/mongodb"
	deliveryGRPC "github.com/Yershuaq/Asik_1_go/inventory-service/internal/delivery/grpc" // Переименовал импорт во избежание конфликта
	"github.com/Yershuaq/Asik_1_go/inventory-service/internal/usecase"
	inventory "github.com/Yershuaq/Asik_1_go/inventory-service/proto" // Предполагается, что proto лежит здесь

	// gRPC и MongoDB импорты...
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

// Функция для периодического обновления кэша
func runCacheRefresher(ctx context.Context, uc *usecase.ProductUseCase, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Refreshing cache...")
			err := uc.InitializeCache(ctx) // Переиспользуем метод инициализации
			if err != nil {
				log.Printf("Error refreshing cache: %v", err)
			}
		case <-ctx.Done():
			log.Println("Stopping cache refresher.")
			return
		}
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI not set")
	}

	// MongoDB Connection
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	// Проверка соединения
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	log.Println("Connected to MongoDB!")

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "inventory" // Default database name
		log.Println("MONGODB_DATABASE not set, using default:", dbName)
	}
	db := client.Database(dbName)

	// Инициализация репозитория
	productRepo := mongodb.NewProductRepository(db)

	// Инициализация кэша
	productCache := cache.NewProductCache()

	// Инициализация use case с репозиторием и кэшем
	productUseCase := usecase.NewProductUseCase(productRepo, productCache)

	// Первоначальная инициализация кэша при старте [cite: 7]
	// Запускаем в отдельной горутине, чтобы не блокировать старт сервера
	go func() {
		// Даем немного времени на установку соединения с БД, если необходимо
		time.Sleep(2 * time.Second)
		err := productUseCase.InitializeCache(context.Background())
		if err != nil {
			log.Printf("Initial cache load failed: %v", err)
		}
	}()

	// Запуск периодического обновления кэша [cite: 7]
	// Используем context для возможности грациозной остановки
	refreshCtx, cancelRefresher := context.WithCancel(context.Background())
	defer cancelRefresher()                                        // Остановит тикер при завершении main
	go runCacheRefresher(refreshCtx, productUseCase, 12*time.Hour) // Обновляем каждые 12 часов

	// gRPC Server Setup
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = ":50051" // Default gRPC port
	}

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	s := grpc.NewServer()
	productServer := deliveryGRPC.NewProductServer(productUseCase) // Используем переименованный импорт
	inventory.RegisterInventoryServiceServer(s, productServer)     // Используем импорт proto

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
