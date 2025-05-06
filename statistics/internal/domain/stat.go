package domain

import "context"

type Stats struct {
	TotalOrders int64
	TotalItems  int64
}

type Repository interface {
	AddInventory(ctx context.Context, itemID string, qty int64) error
	AddOrder(ctx context.Context, userID string, total int64) error
	GetStats(ctx context.Context, userID string) (Stats, error)
}
