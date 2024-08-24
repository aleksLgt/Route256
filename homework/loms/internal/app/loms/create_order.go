package loms

import (
	"context"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"

	"route256/loms/internal/domain"
	servicepb "route256/loms/pkg/api/loms/v1"
	"route256/loms/pkg/prometheus"
)

func (s *Service) CreateOrder(ctx context.Context, in *servicepb.CreateOrderRequest) (*servicepb.CreateOrderResponse, error) {
	handlerName := "POST /v1/order/create"

	ctx, _ = getCtxByTraceID(ctx)

	ctx, span := otel.Tracer("loms").Start(ctx, "handler_create_order")
	defer span.End()

	defer func(createdAt time.Time) {
		prometheus.ObserveGRPCRequestsDurationHistogram(createdAt, "create_order")
	}(time.Now())

	prometheus.IncGRPCRequestsTotalCounter("create_order")

	orderID, err := s.impl.CreateOrder(ctx, in.User, repackItems(in.Items))
	if err != nil {
		return nil, GetErrorResponse(ctx, codes.FailedPrecondition, handlerName, err)
	}

	prometheus.IncGRPCResponseStatusTotalCounter(strconv.FormatUint(uint64(codes.OK), 10), handlerName)

	return &servicepb.CreateOrderResponse{
		OrderID: uint64(*orderID),
	}, nil
}

func repackItems(requestItems []*servicepb.Item) []domain.Item {
	items := make([]domain.Item, len(requestItems))
	for i, item := range requestItems {
		items[i] = domain.Item{
			SKU:   item.Sku,
			Count: item.Count,
		}
	}

	return items
}
