package loms

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"

	"route256/loms/internal/app/definitions/params"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/memory/orders"
	servicepb "route256/loms/pkg/api/loms/v1"
	"route256/loms/pkg/prometheus"
)

func (s *Service) InfoOrder(ctx context.Context, in *servicepb.InfoOrderRequest) (*servicepb.InfoOrderResponse, error) {
	handlerName := fmt.Sprintf("GET /v1/order/{%s}/info", params.ParamOrderID)

	ctx, _ = getCtxByTraceID(ctx)

	ctx, span := otel.Tracer("loms").Start(ctx, "handler_info_order")
	defer span.End()

	defer func(createdAt time.Time) {
		prometheus.ObserveGRPCRequestsDurationHistogram(createdAt, "info_order")
	}(time.Now())

	prometheus.IncGRPCRequestsTotalCounter("info_order")

	order, err := s.impl.InfoOrder(ctx, in.OrderId)
	if err != nil {
		if errors.Is(err, orders.OrderNotFoundError{}) {
			return nil, GetErrorResponse(ctx, codes.NotFound, handlerName, err)
		}

		return nil, GetErrorResponse(ctx, codes.Internal, handlerName, err)
	}

	prometheus.IncGRPCResponseStatusTotalCounter(strconv.FormatUint(uint64(codes.OK), 10), handlerName)

	return repackOrderToProto(order), nil
}

func repackOrderToProto(order *domain.Order) *servicepb.InfoOrderResponse {
	items := make([]*servicepb.Item, len(order.Items))

	for i, n := range order.Items {
		items[i] = &servicepb.Item{
			Sku:   n.SKU,
			Count: n.Count,
		}
	}

	return &servicepb.InfoOrderResponse{
		Status: order.Status,
		User:   order.UserID,
		Items:  items,
	}
}
