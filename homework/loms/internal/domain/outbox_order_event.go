package domain

import "time"

type OutboxOrderEvent struct {
	ID        int64
	OrderID   int64
	EventType EventType
	CreatedAt time.Time
}
