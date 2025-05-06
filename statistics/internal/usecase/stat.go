package usecase

import (
	"context"

	"github.com/Yershuaq/Asik_1_go/statistics/internal/domain"
	events "github.com/Yershuaq/Asik_1_go/statistics/proto/events"
)

type Usecase interface {
	HandleInventory(ctx context.Context, e *events.InventoryEvent) error
	HandleOrder(ctx context.Context, e *events.OrderEvent) error
	GetStats(ctx context.Context, userID string) (domain.Stats, error)
}

type statUC struct {
	repo domain.Repository
}

func New(repo domain.Repository) Usecase {
	return &statUC{repo: repo}
}

func (u *statUC) HandleInventory(ctx context.Context, e *events.InventoryEvent) error {
	return u.repo.AddInventory(ctx, e.ItemId, e.Qty)
}

func (u *statUC) HandleOrder(ctx context.Context, e *events.OrderEvent) error {
	return u.repo.AddOrder(ctx, e.UserId, e.Total)
}

func (u *statUC) GetStats(ctx context.Context, userID string) (domain.Stats, error) {
	return u.repo.GetStats(ctx, userID)
}
