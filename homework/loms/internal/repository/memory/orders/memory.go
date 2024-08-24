package orders

import (
	"context"
	"sync"

	orderStatus "route256/loms/internal/app/definitions"
	"route256/loms/internal/domain"
)

type (
	MemoryStorage struct {
		orders map[int64]domain.Order
		mtx    sync.RWMutex
	}

	OrderNotFoundError struct{}
)

func (_ OrderNotFoundError) Error() string {
	return "Order not found"
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		orders: make(map[int64]domain.Order),
		mtx:    sync.RWMutex{},
	}
}

func (m *MemoryStorage) Create(_ context.Context, userID int64, items []domain.Item) int64 {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	newOrderID := int64(len(m.orders)) + 1

	m.orders[newOrderID] = domain.Order{
		Status: orderStatus.New,
		UserID: userID,
		Items:  items,
	}

	return newOrderID
}

func (m *MemoryStorage) SetStatus(_ context.Context, orderID int64, status string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	order := m.orders[orderID]
	order.Status = status
	m.orders[orderID] = order
}

func (m *MemoryStorage) GetByID(_ context.Context, orderID int64) (*domain.Order, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	order, ok := m.orders[orderID]
	if !ok {
		return nil, OrderNotFoundError{}
	}

	return &order, nil
}
