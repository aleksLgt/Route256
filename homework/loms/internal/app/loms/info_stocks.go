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
	"route256/loms/internal/repository/memory/stocks"
	servicepb "route256/loms/pkg/api/loms/v1"
	"route256/loms/pkg/prometheus"
)

func (s *Service) InfoStocks(ctx context.Context, in *servicepb.InfoStocksRequest) (*servicepb.InfoStocksResponse, error) {
	handlerName := fmt.Sprintf("GET /v1/stock/{%s}/info", params.ParamSKU)

	ctx, _ = getCtxByTraceID(ctx)

	ctx, span := otel.Tracer("loms").Start(ctx, "handler_info_stocks")
	defer span.End()

	defer func(createdAt time.Time) {
		prometheus.ObserveGRPCRequestsDurationHistogram(createdAt, "info_stocks")
	}(time.Now())

	prometheus.IncGRPCRequestsTotalCounter("info_stocks")

	count, err := s.impl.InfoStocks(ctx, in.Sku)
	if err != nil {
		if errors.Is(err, stocks.StockNotFoundError{}) {
			return nil, GetErrorResponse(ctx, codes.NotFound, handlerName, err)
		}

		return nil, GetErrorResponse(ctx, codes.Internal, handlerName, err)
	}

	prometheus.IncGRPCResponseStatusTotalCounter(strconv.FormatUint(uint64(codes.OK), 10), handlerName)

	return &servicepb.InfoStocksResponse{Count: *count}, nil
}
