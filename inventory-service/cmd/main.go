package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	"github.com/Yershuaq/ecommerce/inventory-service/internal/adapter/mongodb"
	"github.com/Yershuaq/ecommerce/inventory-service/internal/delivery/grpc"
	"github.com/Yershuaq/ecommerce/inventory-service/internal/usecase"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Connect(nil); err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(nil)

	db := client.Database("inventory")
	productRepo := mongodb.NewProductRepository(db)
	productUseCase := usecase.NewProductUseCase(productRepo)
	productServer := grpc.NewProductServer(productUseCase)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	proto.RegisterInventoryServiceServer(s, productServer)

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
