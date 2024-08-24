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

	desc "route256/cart/pkg/api/loms/v1"
	"route256/cart/pkg/prometheus"
)

var ErrGetStockInfo = errors.New("LOMSService.InfoStocks failed: ")

func (c *Client) InfoStocks(ctx context.Context, sku int64) (int, error) {
	ctx, span := otel.Tracer("cart").Start(ctx, "loms_client_info_stocks")
	defer span.End()

	defer func(createdAt time.Time) {
		prometheus.ObserveExternalRequestsDurationHistogram(createdAt, "loms", "info_stocks")
	}(time.Now())

	client := desc.NewLOMSClient(c.conn)
	ctx, cancel := context.WithTimeout(ctx, time.Second)

	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "x-auth", c.header)

	prometheus.IncExternalRequestsTotalCounter("loms", "info_stocks")

	traceId := span.SpanContext().TraceID().String()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceId)

	response, err := client.InfoStocks(ctx, &desc.InfoStocksRequest{Sku: uint32(sku)})

	if err != nil {
		prometheus.IncExternalResponseStatusTotalCounter("GET /v1/stock/{sku}/info", strconv.FormatUint(uint64(status.Code(err)), 10))

		return 0, fmt.Errorf("error when calling InfoStocks: %w", err)
	}

	prometheus.IncExternalResponseStatusTotalCounter("GET /v1/stock/{sku}/info", strconv.FormatUint(uint64(codes.OK), 10))

	return int(response.Count), nil
}
