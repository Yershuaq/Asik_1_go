package mongodb

import (
	"context"

	"github.com/Yershuaq/ecommerce/inventory-service/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{
		collection: db.Collection("products"),
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *entity.Product) error {
	product.ID = primitive.NewObjectID().Hex()
	_, err := r.collection.InsertOne(ctx, product)
	return err
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*entity.Product, error) {
	var product entity.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *entity.Product) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": product.ID}, product)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *ProductRepository) List(ctx context.Context, page, limit int32) ([]*entity.Product, int32, error) {
	skip := (page - 1) * limit
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var products []*entity.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return products, int32(total), nil
}
