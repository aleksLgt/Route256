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
	"route256/loms/internal/repository/memory/orders"
	servicepb "route256/loms/pkg/api/loms/v1"
	"route256/loms/pkg/prometheus"
)

func (s *Service) CancelOrder(ctx context.Context, in *servicepb.CancelOrderRequest) (*servicepb.CancelOrderResponse, error) {
	handlerName := fmt.Sprintf("POST /v1/order/{%s}/cancel", params.ParamOrderID)

	ctx, _ = getCtxByTraceID(ctx)

	ctx, span := otel.Tracer("loms").Start(ctx, "handler_cancel_order")
	defer span.End()

	defer func(createdAt time.Time) {
		prometheus.ObserveGRPCRequestsDurationHistogram(createdAt, "cancel_order")
	}(time.Now())

	prometheus.IncGRPCRequestsTotalCounter("cancel_order")

	err := s.impl.CancelOrder(ctx, in.OrderId)
	if err != nil {
		if errors.Is(err, orders.OrderNotFoundError{}) {
			return nil, GetErrorResponse(ctx, codes.NotFound, handlerName, err)
		}

		return nil, GetErrorResponse(ctx, codes.Internal, handlerName, err)
	}

	prometheus.IncGRPCResponseStatusTotalCounter(strconv.FormatUint(uint64(codes.OK), 10), handlerName)

	return &servicepb.CancelOrderResponse{Success: true}, nil
}
