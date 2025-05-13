package usecase

import (
	"context"
	"github.com/Yershuaq/Asik_1_go/inventory-service/internal/adapter/cache" // Импортируем адаптер кэша
	"github.com/Yershuaq/Asik_1_go/inventory-service/internal/entity"
	"github.com/Yershuaq/Asik_1_go/inventory-service/internal/repository" // Предполагается, что у вас есть интерфейс репозитория
	"log"
)

type ProductUseCase struct {
	productRepo repository.ProductRepository // Интерфейс репозитория
	cache       *cache.ProductCache          // Адаптер кэша
}

// Обновляем конструктор
func NewProductUseCase(productRepo repository.ProductRepository, cache *cache.ProductCache) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
		cache:       cache,
	}
}

// CreateProduct: после создания в БД, добавляем в кэш [cite: 5]
func (uc *ProductUseCase) CreateProduct(ctx context.Context, product *entity.Product) error {
	err := uc.productRepo.Create(ctx, product) // Используем интерфейс репозитория
	if err == nil {
		uc.cache.SetProduct(product) // Сохраняем в кэш
		uc.cache.ClearProductLists() // Инвалидируем списки
		log.Printf("Product %s created and added to cache", product.ID)
	}
	return err
}

// GetProductByID: сначала ищем в кэше, потом в БД [cite: 6]
func (uc *ProductUseCase) GetProductByID(ctx context.Context, id string) (*entity.Product, error) {
	// 1. Проверяем кэш
	if product, found := uc.cache.GetProduct(id); found {
		log.Printf("Product %s found in cache", id)
		return product, nil
	}

	// 2. Если в кэше нет, идем в репозиторий (БД)
	log.Printf("Product %s not in cache, fetching from DB", id)
	product, err := uc.productRepo.GetByID(ctx, id) // Используем интерфейс репозитория
	if err != nil {
		return nil, err
	}

	// 3. Сохраняем результат в кэш перед возвратом
	if product != nil {
		uc.cache.SetProduct(product)
	}
	return product, nil
}

// UpdateProduct: обновляем в БД, затем обновляем/инвалидируем кэш
func (uc *ProductUseCase) UpdateProduct(ctx context.Context, product *entity.Product) error {
	err := uc.productRepo.Update(ctx, product) // Используем интерфейс репозитория
	if err == nil {
		uc.cache.SetProduct(product) // Обновляем в кэше
		uc.cache.ClearProductLists() // Инвалидируем списки
		log.Printf("Product %s updated in DB and cache", product.ID)
	}
	return err
}

// DeleteProduct: удаляем из БД, затем удаляем из кэша
func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id string) error {
	err := uc.productRepo.Delete(ctx, id) // Используем интерфейс репозитория
	if err == nil {
		uc.cache.DeleteProduct(id) // Удаляем из кэша
		// Инвалидация списков уже происходит внутри DeleteProduct в адаптере кэша
		log.Printf("Product %s deleted from DB and cache", id)
	}
	return err
}

// ListProducts: сначала ищем в кэше, потом в БД [cite: 6]
func (uc *ProductUseCase) ListProducts(ctx context.Context, page, limit int32) ([]*entity.Product, int32, error) {
	// 1. Проверяем кэш для конкретной страницы
	if products, found := uc.cache.GetProductList(page, limit); found {
		log.Printf("Product list (page %d, limit %d) found in cache", page, limit)
		// Важно: total нужно все равно запросить из БД или кэшировать отдельно
		// Здесь для простоты запрашиваем total из репозитория
		_, total, err := uc.productRepo.List(ctx, page, limit) // Используем интерфейс репозитория
		if err != nil {
			// Можно вернуть кэшированные данные с ошибкой получения total или обработать иначе
			return products, 0, err
		}
		return products, total, nil
	}

	// 2. Если в кэше нет, идем в репозиторий
	log.Printf("Product list (page %d, limit %d) not in cache, fetching from DB", page, limit)
	products, total, err := uc.productRepo.List(ctx, page, limit) // Используем интерфейс репозитория
	if err != nil {
		return nil, 0, err
	}

	// 3. Сохраняем результат в кэш
	if products != nil {
		uc.cache.SetProductList(page, limit, products)
	}
	return products, total, err
}

// InitializeCache: Метод для первоначальной загрузки и периодического обновления кэша [cite: 7]
func (uc *ProductUseCase) InitializeCache(ctx context.Context) error {
	log.Println("Initializing cache...")
	// Получаем все продукты из репозитория (может потребоваться метод ListAll в репозитории)
	// Ограничьте количество или используйте пагинацию, если товаров очень много
	products, _, err := uc.productRepo.List(ctx, 1, 10000) // Пример: загрузка первых 10000
	if err != nil {
		log.Printf("Error fetching products for cache initialization: %v", err)
		return err
	}
	uc.cache.LoadAllProducts(products)
	return nil
}
