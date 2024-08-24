package lomsusecase

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"

	orderStatus "route256/loms/internal/app/definitions"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/memory/stocks"
)

type CreateOrderError struct{}

func (_ CreateOrderError) Error() string {
	return "Error by creating order: "
}

func (s *Service) CreateOrder(ctx context.Context, userID int64, items []domain.Item) (*int64, error) {
	ctx, span := otel.Tracer("loms").Start(ctx, "service_create_order")
	defer span.End()

	orderID, err := s.ordersRepo.Create(ctx, userID, items)

	if err != nil {
		return nil, fmt.Errorf("%w, %w", CreateOrderError{}, err)
	}

	err = s.stocksRepo.Reserve(ctx, items)

	if err != nil {
		err = s.ordersRepo.SetStatus(ctx, orderID, orderStatus.Failed)
		if err != nil {
			return nil, fmt.Errorf("%w, %w", CreateOrderError{}, err)
		}

		return nil, stocks.StockNotFoundError{}
	}

	err = s.ordersRepo.SetStatus(ctx, orderID, orderStatus.AwaitingPayment)
	if err != nil {
		return nil, fmt.Errorf("%w, %w", CreateOrderError{}, err)
	}

	return &orderID, nil
}
