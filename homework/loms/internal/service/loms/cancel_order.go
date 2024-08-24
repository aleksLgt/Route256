package lomsusecase

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"

	orderStatus "route256/loms/internal/app/definitions"
	"route256/loms/internal/repository/memory/orders"
)

type CancelOrderError struct{}

func (_ CancelOrderError) Error() string {
	return "Error by cancelling order: "
}

func (s *Service) CancelOrder(ctx context.Context, orderID int64) error {
	ctx, span := otel.Tracer("loms").Start(ctx, "service_cancel_order")
	defer span.End()

	order, err := s.ordersRepo.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, orders.OrderNotFoundError{}) {
			return orders.OrderNotFoundError{}
		}

		return fmt.Errorf("%w, %w", CancelOrderError{}, err)
	}

	err = s.stocksRepo.ReserveCancel(ctx, order.Items)
	if err != nil {
		return fmt.Errorf("%w, %w", CancelOrderError{}, err)
	}

	err = s.ordersRepo.SetStatus(ctx, orderID, orderStatus.Cancelled)
	if err != nil {
		return fmt.Errorf("%w, %w", CancelOrderError{}, err)
	}

	return nil
}
