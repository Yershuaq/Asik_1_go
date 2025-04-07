package entity

import (
	"time"
)

// OrderItem представляет собой элемент заказа
type OrderItem struct {
	ProductID string  `json:"product_id" bson:"product_id"`
	Quantity  int     `json:"quantity" bson:"quantity"`
	Price     float64 `json:"price" bson:"price"`
}

// Order представляет собой сущность заказа
type Order struct {
	ID        string      `json:"id" bson:"_id,omitempty"`
	UserID    string      `json:"user_id" bson:"user_id"`
	Items     []OrderItem `json:"items" bson:"items"`
	Status    string      `json:"status" bson:"status"`
	Total     float64     `json:"total" bson:"total"`
	CreatedAt time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" bson:"updated_at"`
}

// OrderRepository определяет интерфейс для работы с хранилищем заказов
type OrderRepository interface {
	Create(order *Order) error
	FindByID(id string) (*Order, error)
	Update(order *Order) error
	Delete(id string) error
	FindAll() ([]*Order, error)
	FindByUserID(userID string) ([]*Order, error)
}
