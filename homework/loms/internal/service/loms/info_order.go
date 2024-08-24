package lomsusecase

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"

	"route256/loms/internal/domain"
	"route256/loms/internal/repository/memory/orders"
)

type InfoOrderError struct{}

func (_ InfoOrderError) Error() string {
	return "Error by getting order info: "
}

func (s *Service) InfoOrder(ctx context.Context, orderID int64) (*domain.Order, error) {
	ctx, span := otel.Tracer("loms").Start(ctx, "service_info_order")
	defer span.End()

	order, err := s.ordersRepo.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, orders.OrderNotFoundError{}) {
			return nil, orders.OrderNotFoundError{}
		}

		return nil, fmt.Errorf("%w, %w", InfoOrderError{}, err)
	}

	return order, nil
}
