package producer

import (
	"time"

	"github.com/IBM/sarama"
)

func PrepareConfig(opts ...Option) *sarama.Config {
	c := sarama.NewConfig()

	// алгоритм выбора партиции по ключу
	c.Producer.Partitioner = sarama.NewHashPartitioner

	// acks параметр
	// acks = -1 (all) - ждем успешной записи на лидер партиции и всех in-sync реплик (настроено в кафка кластере)
	c.Producer.RequiredAcks = sarama.WaitForAll

	// семантика exactly once
	c.Producer.Idempotent = false

	// повторы ошибочных отправлений
	// число попыток отправить сообщение
	c.Producer.Retry.Max = 100
	// интервалы между попытками отправить сообщение
	c.Producer.Retry.Backoff = 5 * time.Millisecond

	// Уменьшаем пропускную способность, тем самым гарантируем строгий порядок отправки сообщений/батчей
	c.Net.MaxOpenRequests = 1

	// сжатие на клиенте
	// Если хотим сжимать, то задаем нужный уровень кодировщику
	c.Producer.CompressionLevel = sarama.CompressionLevelDefault
	// И сам кодировщик
	c.Producer.Compression = sarama.CompressionGZIP

	/*
		Если эта конфигурация используется для создания `SyncProducer`, оба параметра должны быть установлены
		в значение true, и вы не не должны читать данные из каналов, поскольку это уже делает продьюсер под капотом.
	*/
	c.Producer.Return.Successes = true
	c.Producer.Return.Errors = true

	for _, opt := range opts {
		_ = opt.Apply(c)
	}

	return c
}
