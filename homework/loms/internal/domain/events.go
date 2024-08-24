package domain

import (
	"errors"
	"time"
)

type EventType string

const (
	EventOrderCreated         EventType = "order-created"
	EventOrderAwaitingPayment EventType = "order-awaiting-payment"
	EventOrderFailed          EventType = "order-failed"
	EventOrderPayed           EventType = "order-payed"
	EventOrderCancelled       EventType = "order-cancelled"
)

type Event struct {
	OrderID         int64     `json:"order_id"`
	ID              int64     `json:"id"`
	EventType       EventType `json:"event"`
	IdempotentKey   string    `json:"idempotent_key"`
	OperationMoment time.Time `json:"moment"`
}

func GetEventTypeByOrderStatus(status string) (EventType, error) {
	switch status {
	case "new":
		return EventOrderCreated, nil
	case "awaiting payment":
		return EventOrderAwaitingPayment, nil
	case "failed":
		return EventOrderFailed, nil
	case "payed":
		return EventOrderPayed, nil
	case "cancelled":
		return EventOrderCancelled, nil
	}

	return "", errors.New("invalid mapping of status to event type")
}
