package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"route256/loms/internal/domain"
	"route256/loms/internal/repository/db/orders"
	"route256/loms/pkg/logger"
	"route256/loms/pkg/producer"
)

type OrdersRepository interface {
	GetUnsentOutboxOrderEvents(ctx context.Context, limit int32) ([]domain.OutboxOrderEvent, error)
	MarkAsSentOutboxOrderEvent(ctx context.Context, eventID int64) error
}

type ProduceOrderEventsJob struct {
	ordersRepository OrdersRepository
	done             chan bool
	connWrite        *pgx.Conn
}

func InitJob(connRead, connWrite *pgx.Conn) *ProduceOrderEventsJob {
	return &ProduceOrderEventsJob{
		ordersRepository: orders.NewStorage(connRead, connWrite),
		done:             make(chan bool),
		connWrite:        connWrite,
	}
}

func (p *ProduceOrderEventsJob) Shutdown() {
	p.done <- true
}

const lomsOrderEventsTopic = "loms.order-events"

var (
	rate  = 3 * time.Second
	limit = 500
)

func (p *ProduceOrderEventsJob) Run() {
	ticker := time.NewTicker(rate)
	defer ticker.Stop()
	defer close(p.done)

	ctx := context.Background()

	for {
		select {
		case <-p.done:
			logger.Infow(ctx, "ProduceOrderEventsJob shutdown complete")
			return
		case <-ticker.C:
			p.processEvents(ctx)
		}
	}
}

func (p *ProduceOrderEventsJob) processEvents(ctx context.Context) {
	logger.Infow(ctx, "ProduceOrderEventsJob processEvents() start")

	events, err := p.ordersRepository.GetUnsentOutboxOrderEvents(ctx, int32(limit))
	failedOrderIds := make(map[int64]struct{})

	if err != nil {
		if fmt.Sprint(err) != "no rows in result set" {
			logger.Errorw(ctx, "Error by getting outbox_order_events", "error", err)
		}

		return
	}

	for _, event := range events {
		if _, ok := failedOrderIds[event.OrderID]; ok {
			// Для сохранения хронологии событий по заказу, новые события не должны быть отправлены в kafka в случае ошибки по более раннему событию
			continue
		}

		err = producer.EmitEvent(lomsOrderEventsTopic, event)
		if err != nil {
			logger.Errorw(ctx, "Error when fixing an event in the kafka queue", "error", err, "orderID", event.OrderID, "event", event.EventType, "topic", lomsOrderEventsTopic)

			failedOrderIds[event.OrderID] = struct{}{}

			continue
		}

		err = p.ordersRepository.MarkAsSentOutboxOrderEvent(ctx, event.ID)
		if err != nil {
			logger.Errorw(ctx, "Error when deleting an outbox order event", "error", err)

			failedOrderIds[event.OrderID] = struct{}{}

			continue
		}
	}

	logger.Infow(ctx, "ProduceOrderEventsJob processEvents() end")
}
