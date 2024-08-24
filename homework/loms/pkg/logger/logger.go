package logger

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerCtx string

const (
	loggerCtxKey loggerCtx = "key"
)

type Logger struct {
	logger *zap.SugaredLogger
}

var global *Logger

func New() (*Logger, error) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level.SetLevel(zapcore.DebugLevel)
	loggerConfig.ErrorOutputPaths = []string{"stdout"}

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("loggerConfig.Build: %w", err)
	}

	once := sync.Once{}
	once.Do(func() {
		global = &Logger{logger: logger.Sugar()}
	})

	return &Logger{logger: logger.Sugar()}, nil
}

func ToContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

func Infow(ctx context.Context, msg string, keysAndValues ...interface{}) {
	// get context
	if loggerC, ok := ctx.Value(loggerCtxKey).(*Logger); ok {
		loggerC.logger.Infow(msg, keysAndValues...)

		return
	}

	if global != nil {
		global.logger.Infow(msg, keysAndValues...)
	}
}

func Errorw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	// get context
	if loggerC, ok := ctx.Value(loggerCtxKey).(*Logger); ok {
		loggerC.logger.Errorw(msg, keysAndValues...)

		return
	}

	if global != nil {
		global.logger.Errorw(msg, keysAndValues...)
	}
}

func Panicw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	// get context
	if loggerC, ok := ctx.Value(loggerCtxKey).(*Logger); ok {
		loggerC.logger.Panicw(msg, keysAndValues...)

		return
	}

	if global != nil {
		global.logger.Panicw(msg, keysAndValues...)
	}
}

func With(args ...interface{}) (*Logger, error) {
	if global == nil {
		return nil, fmt.Errorf("logger not initialized")
	}

	return &Logger{global.logger.With(args...)}, nil
}
