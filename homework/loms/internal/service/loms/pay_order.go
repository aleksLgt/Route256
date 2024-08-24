package lomsusecase

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"

	orderStatus "route256/loms/internal/app/definitions"
	"route256/loms/internal/repository/memory/orders"
)

type PayOrderError struct{}

func (_ PayOrderError) Error() string {
	return "Error by paying order: "
}

func (s *Service) PayOrder(ctx context.Context, orderID int64) error {
	ctx, span := otel.Tracer("loms").Start(ctx, "service_pay_order")
	defer span.End()

	order, err := s.ordersRepo.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, orders.OrderNotFoundError{}) {
			return orders.OrderNotFoundError{}
		}

		return fmt.Errorf("%w, %w", PayOrderError{}, err)
	}

	err = s.stocksRepo.ReserveRemove(ctx, order.Items)
	if err != nil {
		return fmt.Errorf("%w, %w", PayOrderError{}, err)
	}

	err = s.ordersRepo.SetStatus(ctx, orderID, orderStatus.Payed)
	if err != nil {
		return fmt.Errorf("%w, %w", PayOrderError{}, err)
	}

	return nil
}
