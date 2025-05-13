package usecase

import (
	"context"
	"time"

	"github.com/Yershuaq/Asik_1_go/order-service/internal/entity"
)

type OrderUseCase struct {
	repo entity.OrderRepository
}

func NewOrderUseCase(repo entity.OrderRepository) *OrderUseCase {
	return &OrderUseCase{repo: repo}
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, order *entity.Order) error {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.Status = "pending"

	var total float64
	for _, item := range order.Items {
		total += item.Price * float64(item.Quantity)
	}
	order.Total = total

	return uc.repo.Create(order)
}

func (uc *OrderUseCase) GetOrderByID(ctx context.Context, id string) (*entity.Order, error) {
	return uc.repo.FindByID(id)
}

func (uc *OrderUseCase) UpdateOrder(ctx context.Context, order *entity.Order) error {
	order.UpdatedAt = time.Now()
	return uc.repo.Update(order)
}

func (uc *OrderUseCase) DeleteOrder(ctx context.Context, id string) error {
	return uc.repo.Delete(id)
}

func (uc *OrderUseCase) GetAllOrders(ctx context.Context) ([]*entity.Order, error) {
	return uc.repo.FindAll()
}

func (uc *OrderUseCase) GetOrdersByUserID(ctx context.Context, userID string) ([]*entity.Order, error) {
	return uc.repo.FindByUserID(userID)
}
