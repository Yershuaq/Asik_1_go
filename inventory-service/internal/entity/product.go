package entity

import (
	"time"
)

type Product struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Price       float64   `json:"price" bson:"price"`
	Stock       int32     `json:"stock" bson:"stock"`
	Category    string    `json:"category" bson:"category"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

type ProductRepository interface {
	Create(product *Product) error
	FindByID(id string) (*Product, error)
	Update(product *Product) error
	Delete(id string) error
	FindAll() ([]*Product, error)
}
