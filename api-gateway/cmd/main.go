package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// gRPC connections
	inventoryConn, err := grpc.Dial(os.Getenv("INVENTORY_SERVICE_ADDR"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to inventory service: %v", err)
	}
	defer inventoryConn.Close()

	orderConn, err := grpc.Dial(os.Getenv("ORDER_SERVICE_ADDR"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to order service: %v", err)
	}
	defer orderConn.Close()

	// Create gRPC clients
	inventoryClient := inventory.NewInventoryServiceClient(inventoryConn)
	orderClient := order.NewOrderServiceClient(orderConn)

	// Create handlers
	inventoryHandler := http.NewInventoryHandler(inventoryClient)
	orderHandler := http.NewOrderHandler(orderClient)

	// Setup router
	r := gin.Default()

	// Inventory routes
	r.POST("/products", inventoryHandler.CreateProduct)
	r.GET("/products/:id", inventoryHandler.GetProductByID)
	r.PUT("/products/:id", inventoryHandler.UpdateProduct)
	r.DELETE("/products/:id", inventoryHandler.DeleteProduct)
	r.GET("/products", inventoryHandler.ListProducts)

	// Order routes
	r.POST("/orders", orderHandler.CreateOrder)
	r.GET("/orders/:id", orderHandler.GetOrderByID)
	r.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)
	r.GET("/orders", orderHandler.ListUserOrders)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
