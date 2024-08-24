package loms

import (
	"context"
	"fmt"
	"strconv"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"route256/loms/internal/domain"
	servicepb "route256/loms/pkg/api/loms/v1"
	"route256/loms/pkg/logger"
	"route256/loms/pkg/prometheus"
)

var _ servicepb.LOMSServer = (*Service)(nil)

type LOMSService interface {
	CancelOrder(ctx context.Context, orderID int64) error
	CreateOrder(ctx context.Context, userID int64, items []domain.Item) (*int64, error)
	InfoOrder(ctx context.Context, orderID int64) (*domain.Order, error)
	InfoStocks(ctx context.Context, sku uint32) (*int64, error)
	PayOrder(ctx context.Context, orderID int64) error
}

type Service struct {
	servicepb.UnimplementedLOMSServer
	impl LOMSService
}

func NewService(impl LOMSService) *Service {
	return &Service{impl: impl}
}

func GetErrorResponse(ctx context.Context, code codes.Code, handlerName string, err error) error {
	prometheus.IncGRPCResponseStatusTotalCounter(strconv.FormatUint(uint64(code), 10), handlerName)

	return status.Errorf(code, err.Error())
}

func getCtxByTraceID(ctx context.Context) (context.Context, error) {
	// Extract TraceID from header
	md, _ := metadata.FromIncomingContext(ctx)

	if len(md["x-trace-id"]) == 0 {
		logger.Errorw(ctx, "no x-trace-id in the incoming metadata")
		return nil, fmt.Errorf("no x-trace-id")
	}

	traceIdString := md["x-trace-id"][0]

	// Convert string to byte array
	traceId, err := trace.TraceIDFromHex(traceIdString)
	if err != nil {
		logger.Errorw(ctx, "unable to get a TraceID from a hex string", "error", err)
		return nil, fmt.Errorf("trace.TraceIDFromHex: %w", err)
	}

	// Creating a span context with a predefined trace-id
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceId,
	})

	// Embedding span config into the context
	ctx = trace.ContextWithSpanContext(ctx, spanContext)

	return ctx, nil
}
