package usecase

import (
	"context"
	"time"

	"github.com/your-username/ecommerce/order-service/internal/entity"
)

// OrderUseCase реализует бизнес-логику работы с заказами
type OrderUseCase struct {
	repo entity.OrderRepository
}

// NewOrderUseCase создает новый экземпляр OrderUseCase
func NewOrderUseCase(repo entity.OrderRepository) *OrderUseCase {
	return &OrderUseCase{repo: repo}
}

// CreateOrder создает новый заказ
func (uc *OrderUseCase) CreateOrder(ctx context.Context, order *entity.Order) error {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.Status = "pending"

	// Расчет общей суммы заказа
	var total float64
	for _, item := range order.Items {
		total += item.Price * float64(item.Quantity)
	}
	order.Total = total

	return uc.repo.Create(order)
}

// GetOrderByID возвращает заказ по ID
func (uc *OrderUseCase) GetOrderByID(ctx context.Context, id string) (*entity.Order, error) {
	return uc.repo.FindByID(id)
}

// UpdateOrder обновляет информацию о заказе
func (uc *OrderUseCase) UpdateOrder(ctx context.Context, order *entity.Order) error {
	order.UpdatedAt = time.Now()
	return uc.repo.Update(order)
}

// DeleteOrder удаляет заказ
func (uc *OrderUseCase) DeleteOrder(ctx context.Context, id string) error {
	return uc.repo.Delete(id)
}

// GetAllOrders возвращает список всех заказов
func (uc *OrderUseCase) GetAllOrders(ctx context.Context) ([]*entity.Order, error) {
	return uc.repo.FindAll()
}

// GetOrdersByUserID возвращает список заказов пользователя
func (uc *OrderUseCase) GetOrdersByUserID(ctx context.Context, userID string) ([]*entity.Order, error) {
	return uc.repo.FindByUserID(userID)
}
