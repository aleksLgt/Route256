package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/fossoreslp/go-uuid-v4"

	"route256/loms/internal/domain"
	"route256/loms/internal/infra/kafka"
	"route256/loms/internal/infra/kafka/producer"
	"route256/loms/pkg/logger"
)

var global sarama.SyncProducer

func New() (sarama.SyncProducer, error) {
	prod, err := producer.NewSyncProducer(kafka.Config{
		Brokers: []string{
			"kafka0:29092",
		},
	},
		producer.WithIdempotent(),
		producer.WithRequiredAcks(sarama.WaitForAll),
		producer.WithMaxOpenRequests(1),
		producer.WithMaxRetries(5),
		producer.WithRetryBackoff(10*time.Millisecond),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize kafka producer: %w", err)
	}

	once := sync.Once{}
	once.Do(func() {
		global = prod
	})

	return prod, nil
}

func EmitEvent(topicName string, event domain.OutboxOrderEvent) error {
	if global == nil {
		return fmt.Errorf("no producer has been initialized")
	}

	idempotentKey, err := uuid.New()
	if err != nil {
		return fmt.Errorf("generating a new UUID failed: %w", err)
	}

	kafkaEvent := domain.Event{
		ID:              event.ID,
		EventType:       event.EventType,
		OrderID:         event.OrderID,
		IdempotentKey:   idempotentKey.String(),
		OperationMoment: event.CreatedAt,
	}

	bytes, err := json.Marshal(kafkaEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topicName,
		Key:   sarama.StringEncoder(strconv.FormatInt(event.OrderID, 10)),
		Value: sarama.ByteEncoder(bytes),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("app-name"),
				Value: []byte("route256-sync-prod"),
			},
		},
		Timestamp: time.Now(),
	}

	partition, offset, err := global.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("could not send message to Kafka for OrderID %d: %w", event.OrderID, err)
	}

	logger.Infow(context.Background(), "The message was successfully sent to Kafka", "key", event.OrderID, "partition", partition, "offset", offset)

	return nil
}
