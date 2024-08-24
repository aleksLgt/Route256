package stocks

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"

	"route256/loms/internal/domain"
	"route256/loms/pkg/prometheus"
)

type (
	Storage struct {
		connRead  *pgx.Conn
		connWrite *pgx.Conn
		cmdRead   *Queries
		cmdWrite  *Queries
	}

	StockData struct {
		SKU        int32 `json:"sku"`
		TotalCount int32 `json:"total_count"`
		Reserved   int32 `json:"reserved"`
	}

	StockNotFoundError struct{}
)

//go:embed stock-data.json
var stocks []byte

func (_ StockNotFoundError) Error() string {
	return "Stock not found"
}

func NewStorage(connRead, connWrite *pgx.Conn) *Storage {
	return &Storage{
		connRead:  connRead,
		connWrite: connWrite,
		cmdRead:   New(connRead),
		cmdWrite:  New(connWrite),
	}
}

func FillStocks(cmd *Queries) error {
	var fileStocks []StockData

	err := json.Unmarshal(stocks, &fileStocks)
	if err != nil {
		return fmt.Errorf("error when decoding data: %w", err)
	}

	for _, stock := range fileStocks {
		err = cmd.CreateStock(context.Background(), CreateStockParams{
			Sku:        stock.SKU,
			TotalCount: stock.TotalCount,
			Reserved:   stock.Reserved,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) Reserve(ctx context.Context, items []domain.Item) error {
	ctx, span := otel.Tracer("loms").Start(ctx, "db_stocks_reserve")
	defer span.End()

	tx, err := s.connWrite.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error when starting transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	for _, item := range items {
		prometheus.IncDBRequestsTotalCounter("update")

		startTime := time.Now()
		err = s.cmdWrite.WithTx(tx).ReserveStock(ctx, ReserveStockParams{
			Reserved: int32(item.Count),
			Sku:      int32(item.SKU),
		})

		if err != nil {
			prometheus.ObserveDBRequestsDurationHistogram(startTime, "update", "error")

			return err
		}

		prometheus.ObserveDBRequestsDurationHistogram(startTime, "update", "success")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error when commiting transaction: %w", err)
	}

	return nil
}

func (s *Storage) ReserveRemove(ctx context.Context, items []domain.Item) error {
	ctx, span := otel.Tracer("loms").Start(ctx, "db_stocks_reserve_remove")
	defer span.End()

	tx, err := s.connWrite.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error when starting transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	ids := make([]int32, len(items))
	for i, item := range items {
		ids[i] = int32(item.ID)
	}

	prometheus.IncDBRequestsTotalCounter("select")

	startTime := time.Now()
	responseItems, err := s.cmdWrite.GetStocks(ctx, ids)

	if err != nil {
		prometheus.ObserveDBRequestsDurationHistogram(startTime, "select", "error")

		return fmt.Errorf("error when getting stocks: %w", err)
	}

	prometheus.ObserveDBRequestsDurationHistogram(startTime, "select", "success")

	for _, responseItem := range responseItems {
		prometheus.IncDBRequestsTotalCounter("update")

		startTime = time.Now()
		err = s.cmdWrite.WithTx(tx).RemoveReserveStock(ctx, RemoveReserveStockParams{
			TotalCount: responseItem.TotalCount - responseItem.Reserved,
			Reserved:   0,
			ID:         responseItem.Sku,
		})

		if err != nil {
			prometheus.ObserveDBRequestsDurationHistogram(startTime, "update", "error")

			return fmt.Errorf("error when removing reserved stock: %w", err)
		}

		prometheus.ObserveDBRequestsDurationHistogram(startTime, "update", "success")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error when commiting transaction: %w", err)
	}

	return nil
}

func (s *Storage) ReserveCancel(ctx context.Context, items []domain.Item) error {
	ctx, span := otel.Tracer("loms").Start(ctx, "db_stocks_reserve_cancel")
	defer span.End()

	tx, err := s.connWrite.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error when starting transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	for _, item := range items {
		prometheus.IncDBRequestsTotalCounter("update")

		startTime := time.Now()
		err = s.cmdWrite.WithTx(tx).ReserveStock(ctx, ReserveStockParams{
			Reserved: 0,
			Sku:      int32(item.SKU),
		})

		if err != nil {
			prometheus.ObserveDBRequestsDurationHistogram(startTime, "update", "error")

			return fmt.Errorf("error when reserving stock: %w", err)
		}

		prometheus.ObserveDBRequestsDurationHistogram(startTime, "update", "success")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error when commiting transaction: %w", err)
	}

	return nil
}

func (s *Storage) GetBySKU(ctx context.Context, sku uint32) (*int64, error) {
	ctx, span := otel.Tracer("loms").Start(ctx, "db_stocks_get_by_sku")
	defer span.End()

	prometheus.IncDBRequestsTotalCounter("select")

	startTime := time.Now()
	stock, err := s.cmdRead.GetStock(ctx, int32(sku))

	if err != nil {
		prometheus.ObserveDBRequestsDurationHistogram(startTime, "select", "error")
		return nil, StockNotFoundError{}
	}

	prometheus.ObserveDBRequestsDurationHistogram(startTime, "select", "success")

	if stock.Reserved > 0 {
		count := int64(0)
		return &count, nil
	}

	stockCount := int64(stock.TotalCount)

	return &stockCount, nil
}
