package memorycartrepo

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"

	"route256/cart/internal/domain"
	"route256/cart/pkg/prometheus"
)

type (
	itemsMap map[int64]domain.Item

	MemoryStorage struct {
		items map[int64]itemsMap
		mtx   sync.RWMutex
	}

	CartItemsNotFoundError struct{}
)

func (_ CartItemsNotFoundError) Error() string {
	return "CartItems not found"
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		items: make(map[int64]itemsMap),
		mtx:   sync.RWMutex{},
	}
}

func (m *MemoryStorage) getTotalCountItems() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	count := 0
	for _, userItems := range m.items {
		count += len(userItems)
	}

	return count
}

func (m *MemoryStorage) Add(ctx context.Context, userID int64, item domain.Item) {
	_, span := otel.Tracer("cart").Start(ctx, "memory_add")
	defer span.End()

	m.mtx.Lock()
	if m.items[userID] == nil {
		m.items[userID] = itemsMap{}
	}

	m.items[userID][item.SKU] = domain.Item{
		SKU:   item.SKU,
		Count: m.items[userID][item.SKU].Count + item.Count,
	}
	m.mtx.Unlock()

	prometheus.UpdateMemoryCartItemsTotalCounter(m.getTotalCountItems())
}

func (m *MemoryStorage) DeleteOne(ctx context.Context, userID, skuID int64) {
	_, span := otel.Tracer("cart").Start(ctx, "memory_delete_one")
	defer span.End()

	m.mtx.RLock()
	userItemsMap, ok := m.items[userID]
	m.mtx.RUnlock()

	if !ok {
		return
	}

	m.mtx.Lock()
	delete(userItemsMap, skuID)
	m.mtx.Unlock()

	prometheus.UpdateMemoryCartItemsTotalCounter(m.getTotalCountItems())
}

func (m *MemoryStorage) DeleteAll(ctx context.Context, userID int64) {
	_, span := otel.Tracer("cart").Start(ctx, "memory_delete_all")
	defer span.End()

	m.mtx.Lock()
	delete(m.items, userID)
	m.mtx.Unlock()

	prometheus.UpdateMemoryCartItemsTotalCounter(m.getTotalCountItems())
}

func (m *MemoryStorage) GetAll(ctx context.Context, userID int64) ([]domain.Item, error) {
	_, span := otel.Tracer("cart").Start(ctx, "memory_get_all")
	defer span.End()

	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if currentSkuItems, ok := m.items[userID]; ok {
		if len(currentSkuItems) == 0 {
			return nil, CartItemsNotFoundError{}
		}

		skuItems := make([]domain.Item, len(currentSkuItems))

		idx := 0

		for _, currentSkuItem := range currentSkuItems {
			skuItems[idx] = currentSkuItem
			idx++
		}

		return skuItems, nil
	}

	return nil, CartItemsNotFoundError{}
}

func (m *MemoryStorage) GetAllOld(_ context.Context, userID int64) ([]domain.Item, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	var skuItems []domain.Item

	if currentSkuItems, ok := m.items[userID]; ok {
		if len(currentSkuItems) == 0 {
			return nil, CartItemsNotFoundError{}
		}

		for _, currentSkuItem := range currentSkuItems {
			skuItems = append(skuItems, currentSkuItem)
		}

		return skuItems, nil
	}

	return nil, CartItemsNotFoundError{}
}
