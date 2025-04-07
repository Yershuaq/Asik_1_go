package usecase

import (
	"context"
	"time"

	"github.com/your-username/ecommerce/inventory-service/internal/entity"
)

// ProductUseCase реализует бизнес-логику работы с товарами
type ProductUseCase struct {
	repo entity.ProductRepository
}

// NewProductUseCase создает новый экземпляр ProductUseCase
func NewProductUseCase(repo entity.ProductRepository) *ProductUseCase {
	return &ProductUseCase{repo: repo}
}

// CreateProduct создает новый товар
func (uc *ProductUseCase) CreateProduct(ctx context.Context, product *entity.Product) error {
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	return uc.repo.Create(product)
}

// GetProductByID возвращает товар по ID
func (uc *ProductUseCase) GetProductByID(ctx context.Context, id string) (*entity.Product, error) {
	return uc.repo.FindByID(id)
}

// UpdateProduct обновляет информацию о товаре
func (uc *ProductUseCase) UpdateProduct(ctx context.Context, product *entity.Product) error {
	product.UpdatedAt = time.Now()
	return uc.repo.Update(product)
}

// DeleteProduct удаляет товар
func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id string) error {
	return uc.repo.Delete(id)
}

// GetAllProducts возвращает список всех товаров
func (uc *ProductUseCase) GetAllProducts(ctx context.Context) ([]*entity.Product, error) {
	return uc.repo.FindAll()
}
