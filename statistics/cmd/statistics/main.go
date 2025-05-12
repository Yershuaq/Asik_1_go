package main

import (
	"log"
	"os"

	"github.com/Yershuaq/Asik_1_go/statistics/internal/adapter/db"
	grpcAdapter "github.com/Yershuaq/Asik_1_go/statistics/internal/adapter/grpc"
	natsAdapter "github.com/Yershuaq/Asik_1_go/statistics/internal/adapter/nats"
	"github.com/Yershuaq/Asik_1_go/statistics/internal/usecase"
	"github.com/nats-io/nats.go"
)

func main() {
	// 1. Подключаемся к NATS
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatal("NATS connect:", err)
	}

	// 2. Подключаемся к MongoDB
	repo, err := db.New(os.Getenv("MONGO_URI"))
	if err != nil {
		log.Fatal("Mongo connect:", err)
	}

	// 3. Инициализируем бизнес‑логику
	uc := usecase.New(repo)

	// 4. Запускаем NATS‑подписчика
	if err := natsAdapter.Start(nc, uc); err != nil {
		log.Fatal("NATS subscribe:", err)
	}

	// 5. Запускаем gRPC‑сервер
	if err := grpcAdapter.Run(uc); err != nil {
		log.Fatal("gRPC serve:", err)
	}
}
