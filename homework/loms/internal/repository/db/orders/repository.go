package orders

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"

	orderStatus "route256/loms/internal/app/definitions"
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

	OrderNotFoundError struct{}
)

func (_ OrderNotFoundError) Error() string {
	return "Order not found"
}

func NewStorage(connRead, connWrite *pgx.Conn) *Storage {
	return &Storage{
		connRead:  connRead,
		connWrite: connWrite,
		cmdRead:   New(connRead),
		cmdWrite:  New(connWrite),
	}
}

func (s *Storage) Create(ctx context.Context, userID int64, items []domain.Item) (int64, error) {
	ctx, span := otel.Tracer("loms").Start(ctx, "db_orders_create")
	defer span.End()

	tx, err := s.connWrite.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("could not begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	prometheus.IncDBRequestsTotalCounter("insert")

	startTime := time.Now()
	orderID, err := s.cmdWrite.WithTx(tx).CreateOrder(ctx, CreateOrderParams{
		UserID: userID,
		Status: orderStatus.New,
	})

	if err != nil {
		prometheus.ObserveDBRequestsDurationHistogram(startTime, "insert", "error")
		return 0, err
	}

	prometheus.ObserveDBRequestsDurationHistogram(startTime, "insert", "success")

	for _, item := range items {
		prometheus.IncDBRequestsTotalCounter("insert")

		startTime = time.Now()
		err = s.cmdWrite.WithTx(tx).CreateOrderItem(ctx, CreateOrderItemParams{
			OrderID: int64(orderID),
			Sku:     int32(item.SKU),
			Count:   int32(item.Count),
		})

		if err != nil {
			prometheus.ObserveDBRequestsDurationHistogram(startTime, "insert", "error")
			return 0, err
		}

		prometheus.ObserveDBRequestsDurationHistogram(startTime, "insert", "success")
	}

	prometheus.IncDBRequestsTotalCounter("insert")

	startTime = time.Now()
	err = s.cmdWrite.WithTx(tx).CreateOutboxOrderEvent(ctx, CreateOutboxOrderEventParams{
		OrderID:   int64(orderID),
		EventType: string(domain.EventOrderCreated),
	})

	if err != nil {
		prometheus.ObserveDBRequestsDurationHistogram(startTime, "insert", "error")
		return 0, err
	}

	prometheus.ObserveDBRequestsDurationHistogram(startTime, "insert", "success")

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("could not commit transaction: %w", err)
	}

	return int64(orderID), nil
}

func (s *Storage) SetStatus(ctx context.Context, orderID int64, status string) error {
	ctx, span := otel.Tracer("loms").Start(ctx, "db_orders_set_status")
	defer span.End()

	tx, err := s.connWrite.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	prometheus.IncDBRequestsTotalCounter("update")

	startTime := time.Now()
	err = s.cmdWrite.WithTx(tx).SetOrderStatus(ctx, SetOrderStatusParams{
		Status: status,
		ID:     int32(orderID),
	})

	if err != nil {
		prometheus.ObserveDBRequestsDurationHistogram(startTime, "update", "error")
		return err
	}

	prometheus.ObserveDBRequestsDurationHistogram(startTime, "update", "success")

	eventType, err := domain.GetEventTypeByOrderStatus(status)
	if err != nil {
		return fmt.Errorf("incorrect event type: %w", err)
	}

	prometheus.IncDBRequestsTotalCounter("insert")

	startTime = time.Now()
	err = s.cmdWrite.WithTx(tx).CreateOutboxOrderEvent(ctx, CreateOutboxOrderEventParams{
		OrderID:   orderID,
		EventType: string(eventType),
	})

	if err != nil {
		prometheus.ObserveDBRequestsDurationHistogram(startTime, "insert", "error")
		return err
	}

	prometheus.ObserveDBRequestsDurationHistogram(startTime, "insert", "success")

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (s *Storage) GetByID(ctx context.Context, orderID int64) (*domain.Order, error) {
	ctx, span := otel.Tracer("loms").Start(ctx, "db_orders_get_by_id")
	defer span.End()

	prometheus.IncDBRequestsTotalCounter("select")

	startTime := time.Now()
	orderResponse, err := s.cmdRead.GetOrder(ctx, int32(orderID))

	if err != nil {
		prometheus.ObserveDBRequestsDurationHistogram(startTime, "select", "error")
		return nil, OrderNotFoundError{}
	}

	prometheus.ObserveDBRequestsDurationHistogram(startTime, "select", "success")

	order := repackOrder(orderResponse)

	prometheus.IncDBRequestsTotalCounter("select")

	startTime = time.Now()
	itemsResponse, err := s.cmdRead.GetOrderItems(ctx, orderID)

	if err != nil {
		prometheus.ObserveDBRequestsDurationHistogram(startTime, "select", "error")
		return nil, err
	}

	prometheus.ObserveDBRequestsDurationHistogram(startTime, "select", "success")

	order.Items = repackItems(itemsResponse)

	return &order, nil
}

func (s *Storage) GetUnsentOutboxOrderEvents(ctx context.Context, limit int32) ([]domain.OutboxOrderEvent, error) {
	ctx, span := otel.Tracer("loms").Start(ctx, "db_orders_get_unsent_outbox_order_event")
	defer span.End()

	prometheus.IncDBRequestsTotalCounter("select")

	startTime := time.Now()
	orderEventsResponse, err := s.cmdWrite.GetUnsentOutboxOrderEvents(ctx, limit)

	if err != nil {
		prometheus.ObserveDBRequestsDurationHistogram(startTime, "select", "error")
		return nil, err
	}

	prometheus.ObserveDBRequestsDurationHistogram(startTime, "select", "success")

	events := repackOutboxOrderEvents(orderEventsResponse)

	return events, nil
}

func (s *Storage) MarkAsSentOutboxOrderEvent(ctx context.Context, eventID int64) error {
	ctx, span := otel.Tracer("loms").Start(ctx, "db_orders_mark_as_sent_outbox_order_event")
	defer span.End()

	prometheus.IncDBRequestsTotalCounter("update")

	startTime := time.Now()
	err := s.cmdWrite.MarkAsSentOutboxOrderEvent(ctx, int32(eventID))

	if err != nil {
		prometheus.ObserveDBRequestsDurationHistogram(startTime, "update", "error")
		return err
	}

	prometheus.ObserveDBRequestsDurationHistogram(startTime, "update", "success")

	return nil
}

func repackOrder(order Order) domain.Order {
	return domain.Order{
		ID:     int64(order.ID),
		UserID: order.UserID,
		Status: order.Status,
	}
}

func repackItems(responseItems []OrderItem) []domain.Item {
	items := make([]domain.Item, len(responseItems))
	for i, responseItem := range responseItems {
		items[i] = domain.Item{
			ID:      int64(responseItem.ID),
			OrderID: responseItem.OrderID,
			SKU:     uint32(responseItem.Sku),
			Count:   uint32(responseItem.Count),
		}
	}

	return items
}

func repackOutboxOrderEvents(events []OutboxOrderEvent) []domain.OutboxOrderEvent {
	responseEvents := make([]domain.OutboxOrderEvent, len(events))
	for i, event := range events {
		responseEvents[i] = domain.OutboxOrderEvent{
			ID:        int64(event.ID),
			OrderID:   event.OrderID,
			EventType: domain.EventType(event.EventType),
			CreatedAt: event.CreatedAt.Time,
		}
	}

	return responseEvents
}
