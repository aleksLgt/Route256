package lomsusecase

import (
	"context"

	"route256/loms/internal/domain"
)

type (
	OrdersRepository interface {
		Create(_ context.Context, userID int64, items []domain.Item) (int64, error)
		SetStatus(_ context.Context, orderID int64, status string) error
		GetByID(_ context.Context, orderID int64) (*domain.Order, error)
	}
	StocksRepository interface {
		Reserve(_ context.Context, items []domain.Item) error
		ReserveRemove(_ context.Context, items []domain.Item) error
		ReserveCancel(_ context.Context, items []domain.Item) error
		GetBySKU(_ context.Context, sku uint32) (*int64, error)
	}

	Service struct {
		ordersRepo OrdersRepository
		stocksRepo StocksRepository
	}
)

func NewService(ordersRepo OrdersRepository, stocksRepo StocksRepository) *Service {
	return &Service{
		ordersRepo: ordersRepo,
		stocksRepo: stocksRepo,
	}
}
