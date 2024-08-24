package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"

	"route256/notifier/internal/infra/kafka"
	"route256/notifier/internal/infra/kafka/consumer_group"
	"route256/notifier/pkg/logger"
)

type flags struct {
	topic             string
	bootstrapServer   string
	consumerGroupName string
}

var cliFlags = flags{}

func newConfig(f flags) kafka.Config {
	return kafka.Config{
		Brokers: []string{
			f.bootstrapServer,
		},
	}
}

func init() {
	flag.StringVar(&cliFlags.topic, "topic", "loms.order-events", "topic to produce")
	flag.StringVar(&cliFlags.bootstrapServer, "bootstrap-server", "localhost:9092", "kafka broker host and port")
	flag.StringVar(&cliFlags.consumerGroupName, "cg-name", "route256-consumer-group", "topic to produce")

	flag.Parse()
}

func main() {
	_, err := logger.New()
	if err != nil {
		panic(err)
	}

	loggerCustom, err := logger.With("service", "notifier")
	if err != nil {
		// panic, as an error can only occur if the logger is not initialized
		panic(err)
	}

	ctx := logger.ToContext(context.Background(), loggerCustom)

	var (
		wg   = &sync.WaitGroup{}
		conf = newConfig(cliFlags)
	)

	ctx = runSignalHandler(ctx, wg)

	handler := consumer_group.NewConsumerGroupHandler()

	cg, err := consumer_group.NewConsumerGroup(
		conf.Brokers,
		cliFlags.consumerGroupName,
		[]string{cliFlags.topic},
		handler,
		consumer_group.WithOffsetsInitial(sarama.OffsetOldest),
	)
	if err != nil {
		logger.Panicw(ctx, "Failed to create consumer group", "err", err)
	}

	defer cg.Close()

	runCGErrorHandler(ctx, cg, wg)

	cg.Run(ctx, wg)

	wg.Wait()
}

func runSignalHandler(ctx context.Context, wg *sync.WaitGroup) context.Context {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	sigCtx, cancel := context.WithCancel(ctx)

	wg.Add(1)

	go func() {
		defer signal.Stop(sigterm)
		defer wg.Done()
		defer cancel()

		for {
			select {
			case sig, ok := <-sigterm:
				if !ok {
					logger.Debugw(ctx, "[signal] signal chan closed", "signal", sig.String())
					return
				}

				logger.Infow(ctx, "[signal] signal received", "signal", sig.String())

				return
			case _, ok := <-sigCtx.Done():
				if !ok {
					logger.Debugw(ctx, "[signal] context closed")
					return
				}

				logger.Errorw(ctx, "[signal] ctx done", "err", ctx.Err())

				return
			}
		}
	}()

	return sigCtx
}

func runCGErrorHandler(ctx context.Context, cg sarama.ConsumerGroup, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case chErr, ok := <-cg.Errors():
				if !ok {
					logger.Debugw(ctx, "[cg-error] chan closed")
					return
				}

				logger.Errorw(ctx, "[cg-error] error", "err", chErr)
			case <-ctx.Done():
				logger.Errorw(ctx, "[cg-error] ctx closed", "err", ctx.Err())
				return
			}
		}
	}()
}
