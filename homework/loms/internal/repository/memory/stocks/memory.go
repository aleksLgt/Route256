package stocks

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"

	"route256/loms/internal/domain"
)

type (
	MemoryStorage struct {
		stocks map[uint32]*domain.Stock
		mtx    sync.RWMutex
	}

	StockNotFoundError struct{}
)

//go:embed stock-data.json
var stocks []byte

func (_ StockNotFoundError) Error() string {
	return "Stock not found"
}

func NewMemoryStorage() (*MemoryStorage, error) {
	var fileStocks []Stock

	err := json.Unmarshal(stocks, &fileStocks)
	if err != nil {
		return nil, fmt.Errorf("error when decoding data: %w", err)
	}

	initStocks := make(map[uint32]*domain.Stock, len(fileStocks))
	for _, stock := range fileStocks {
		initStocks[stock.SKU] = &domain.Stock{
			Sku:        stock.SKU,
			TotalCount: stock.TotalCount,
			Reserved:   stock.Reserved,
		}
	}

	return &MemoryStorage{
		stocks: initStocks,
		mtx:    sync.RWMutex{},
	}, nil
}

type Stock struct {
	SKU        uint32 `json:"sku"`
	TotalCount int64  `json:"total_count"`
	Reserved   int64  `json:"reserved"`
}

func (m *MemoryStorage) Reserve(_ context.Context, items []domain.Item) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for _, item := range items {
		if _, ok := m.stocks[item.SKU]; !ok {
			return StockNotFoundError{}
		}

		m.stocks[item.SKU].Reserved = int64(item.Count)
	}

	return nil
}

func (m *MemoryStorage) ReserveRemove(_ context.Context, items []domain.Item) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for _, item := range items {
		m.stocks[item.SKU].TotalCount -= m.stocks[item.SKU].Reserved
		m.stocks[item.SKU].Reserved = 0
	}
}

func (m *MemoryStorage) ReserveCancel(_ context.Context, items []domain.Item) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for _, item := range items {
		m.stocks[item.SKU].Reserved = 0
	}
}

func (m *MemoryStorage) GetBySKU(_ context.Context, sku uint32) (*int64, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	stock, ok := m.stocks[sku]
	if !ok {
		return nil, StockNotFoundError{}
	}

	if stock.Reserved > 0 {
		count := int64(0)
		return &count, nil
	}

	return &stock.TotalCount, nil
}
