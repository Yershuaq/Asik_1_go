package usecase

import (
	"context"

	"github.com/Yershuaq/ecommerce/inventory-service/internal/entity"
	"github.com/Yershuaq/ecommerce/inventory-service/internal/repository"
)

type ProductUseCase struct {
	productRepo repository.ProductRepository
}

func NewProductUseCase(productRepo repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
	}
}

func (uc *ProductUseCase) CreateProduct(ctx context.Context, product *entity.Product) error {
	return uc.productRepo.Create(ctx, product)
}

func (uc *ProductUseCase) GetProductByID(ctx context.Context, id string) (*entity.Product, error) {
	return uc.productRepo.GetByID(ctx, id)
}

func (uc *ProductUseCase) UpdateProduct(ctx context.Context, product *entity.Product) error {
	return uc.productRepo.Update(ctx, product)
}

func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id string) error {
	return uc.productRepo.Delete(ctx, id)
}

func (uc *ProductUseCase) ListProducts(ctx context.Context, page, limit int32) ([]*entity.Product, int32, error) {
	return uc.productRepo.List(ctx, page, limit)
}
