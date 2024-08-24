package loms

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"route256/cart/internal/domain"
	desc "route256/cart/pkg/api/loms/v1"
	"route256/cart/pkg/prometheus"
)

var ErrCreateOrder = errors.New("LOMSService.CreateOrder failed: ")

func (c *Client) CreateOrder(ctx context.Context, userID int64, items []domain.Item) (int, error) {
	ctx, span := otel.Tracer("cart").Start(ctx, "loms_client_create_order")
	defer span.End()

	defer func(createdAt time.Time) {
		prometheus.ObserveExternalRequestsDurationHistogram(createdAt, "loms", "create_order")
	}(time.Now())

	client := desc.NewLOMSClient(c.conn)
	ctx, cancel := context.WithTimeout(ctx, time.Second)

	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "x-auth", c.header)

	responseItems := repackItems(items)

	prometheus.IncExternalRequestsTotalCounter("loms", "create_order")

	traceId := span.SpanContext().TraceID().String()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceId)

	response, err := client.CreateOrder(ctx, &desc.CreateOrderRequest{
		User:  userID,
		Items: responseItems,
	})

	if err != nil {
		prometheus.IncExternalResponseStatusTotalCounter("POST /v1/order/create", strconv.FormatUint(uint64(status.Code(err)), 10))

		return 0, fmt.Errorf("error when calling CreateOrder: %w", err)
	}

	prometheus.IncExternalResponseStatusTotalCounter("POST /v1/order/create", strconv.FormatUint(uint64(codes.OK), 10))

	return int(response.OrderID), nil
}

func repackItems(requestItems []domain.Item) []*desc.Item {
	items := make([]*desc.Item, len(requestItems))
	for i, item := range requestItems {
		items[i] = &desc.Item{
			Sku:   uint32(item.SKU),
			Count: uint32(item.Count),
		}
	}

	return items
}
