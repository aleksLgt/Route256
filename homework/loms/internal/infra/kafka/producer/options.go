package producer

import (
	"time"

	"github.com/IBM/sarama"
)

// Option is a configuration callback.
type Option interface {
	Apply(*sarama.Config) error
}

type optionFn func(*sarama.Config) error

func (fn optionFn) Apply(c *sarama.Config) error {
	return fn(c)
}

func WithRequiredAcks(acks sarama.RequiredAcks) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Producer.RequiredAcks = acks
		return nil
	})
}

func WithIdempotent() Option {
	return optionFn(func(c *sarama.Config) error {
		c.Producer.Idempotent = true
		return nil
	})
}

func WithMaxRetries(n int) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Producer.Retry.Max = n
		return nil
	})
}

func WithRetryBackoff(d time.Duration) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Producer.Retry.Backoff = d
		return nil
	})
}

func WithMaxOpenRequests(n int) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Net.MaxOpenRequests = n
		return nil
	})
}
