package lomsusecase

import (
	"context"

	"go.opentelemetry.io/otel"

	"route256/loms/internal/repository/memory/stocks"
)

func (s *Service) InfoStocks(ctx context.Context, sku uint32) (*int64, error) {
	ctx, span := otel.Tracer("loms").Start(ctx, "service_info_stocks")
	defer span.End()

	count, err := s.stocksRepo.GetBySKU(ctx, sku)
	if err != nil {
		return nil, stocks.StockNotFoundError{}
	}

	return count, nil
}
