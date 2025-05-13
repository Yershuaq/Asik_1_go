package cache

import (
	"fmt"
	"github.com/Yershuaq/Asik_1_go/inventory-service/internal/entity"
	"github.com/patrickmn/go-cache"
	"log"
	"strings"
	"time"
)

// Константы для ключей кэша и времени жизни
const (
	productListKeyPrefix    = "products_list_"
	productDetailsKeyPrefix = "product_details_"
	defaultExpiration       = 12 * time.Hour // Требование 1.e [cite: 7]
	cleanupInterval         = 1 * time.Hour
)

type ProductCache struct {
	cache *cache.Cache
}

func NewProductCache() *ProductCache {
	// Инициализация кэша с временем жизни по умолчанию и интервалом очистки
	c := cache.New(defaultExpiration, cleanupInterval)
	return &ProductCache{cache: c}
}

// GetProduct получает товар из кэша по ID
func (pc *ProductCache) GetProduct(id string) (*entity.Product, bool) {
	key := productDetailsKeyPrefix + id
	if product, found := pc.cache.Get(key); found {
		return product.(*entity.Product), true
	}
	return nil, false
}

// SetProduct сохраняет товар в кэш
func (pc *ProductCache) SetProduct(product *entity.Product) {
	key := productDetailsKeyPrefix + product.ID
	pc.cache.Set(key, product, cache.DefaultExpiration) // Используем время жизни по умолчанию (12 часов)
}

// DeleteProduct удаляет товар из кэша
func (pc *ProductCache) DeleteProduct(id string) {
	key := productDetailsKeyPrefix + id
	pc.cache.Delete(key)
	// Также нужно инвалидировать кэш списков, если он используется
	// Например, можно удалить все списки или использовать более сложную логику
	pc.ClearProductLists() // Пример простой инвалидации
}

// --- Кэширование списков товаров (Пример) ---

// GetProductList получает список товаров из кэша
// Ключ может зависеть от параметров пагинации (page, limit)
func (pc *ProductCache) GetProductList(page, limit int32) ([]*entity.Product, bool) {
	key := fmt.Sprintf("%s%d_%d", productListKeyPrefix, page, limit) // Формируем ключ для конкретной страницы
	if products, found := pc.cache.Get(key); found {
		return products.([]*entity.Product), true
	}
	return nil, false
}

// SetProductList сохраняет список товаров в кэш
func (pc *ProductCache) SetProductList(page, limit int32, products []*entity.Product) {
	key := fmt.Sprintf("%s%d_%d", productListKeyPrefix, page, limit)
	pc.cache.Set(key, products, cache.DefaultExpiration)
}

// ClearProductLists очищает кэш списков товаров (для инвалидации)
func (pc *ProductCache) ClearProductLists() {
	// Простой вариант: перебрать ключи и удалить те, что начинаются с productListKeyPrefix
	items := pc.cache.Items()
	for k := range items {
		if strings.HasPrefix(k, productListKeyPrefix) {
			pc.cache.Delete(k)
		}
	}
}

// LoadAllProducts загружает все товары (например, из репозитория) и кэширует их
// Вызывается при старте сервиса и периодически для обновления [cite: 7]
func (pc *ProductCache) LoadAllProducts(products []*entity.Product) {
	// Можно очистить старый кэш перед загрузкой
	pc.cache.Flush() // Очистка всего кэша
	for _, p := range products {
		pc.SetProduct(p) // Кэшируем каждый товар отдельно
	}
	// Здесь можно также предзаполнить кэш для популярных списков, если нужно
	log.Println("Cache initialized/refreshed with", len(products), "products")
}
