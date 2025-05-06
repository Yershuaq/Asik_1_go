package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Yershuaq/Asik_1_go/statistics/internal/domain"
)

type Repo struct {
	client *mongo.Client
}

func New(uri string) (domain.Repository, error) {
	cli, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &Repo{client: cli}, nil
}

func (r *Repo) AddInventory(ctx context.Context, itemID string, qty int64) error {
	_, err := r.client.Database("stats").
		Collection("inventoryStats").
		UpdateOne(ctx,
			bson.M{"item_id": itemID},
			bson.M{"$inc": bson.M{"qty": qty}},
			options.Update().SetUpsert(true),
		)
	return err
}

func (r *Repo) AddOrder(ctx context.Context, userID string, total int64) error {
	_, err := r.client.Database("stats").
		Collection("orderStats").
		UpdateOne(ctx,
			bson.M{"user_id": userID},
			bson.M{"$inc": bson.M{"total": total}},
			options.Update().SetUpsert(true),
		)
	return err
}

func (r *Repo) GetStats(ctx context.Context, userID string) (domain.Stats, error) {
	var inv struct{ Qty int64 }
	var ord struct{ Total int64 }
	_ = r.client.Database("stats").Collection("inventoryStats").FindOne(ctx, bson.M{"item_id": userID}).Decode(&inv)
	_ = r.client.Database("stats").Collection("orderStats").FindOne(ctx, bson.M{"user_id": userID}).Decode(&ord)
	return domain.Stats{TotalItems: inv.Qty, TotalOrders: ord.Total}, nil
}
