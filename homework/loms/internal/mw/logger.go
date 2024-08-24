package mw

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"route256/loms/pkg/logger"
)

func Logger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	ctx, span := otel.Tracer("loms").Start(ctx, "loms_middleware_logger")
	defer span.End()

	traceId := span.SpanContext().TraceID().String()
	spanId := span.SpanContext().SpanID().String()

	raw, _ := protojson.Marshal((req).(proto.Message))
	logger.Infow(ctx, "request", "service", "loms", "method", info.FullMethod, "req", string(raw), "trace_id", traceId, "span_id", spanId)

	resp, err := handler(ctx, req)
	if err != nil {
		logger.Errorw(ctx, "response", "service", "loms", "method", info.FullMethod, "err", err, "trace_id", traceId, "span_id", spanId)
		return nil, err
	}

	rawResp, _ := protojson.Marshal((resp).(proto.Message))

	logger.Infow(ctx, "response", "service", "loms", "method", info.FullMethod, "resp", string(rawResp), "trace_id", traceId, "span_id", spanId)

	return resp, nil
}

func WithHTTPLoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		logger.Infow(context.Background(), "incoming request",
			"method", r.Method,
			"host", r.Host,
			"path", r.URL.Path,
			"user_agent", r.UserAgent(),
			"proto", r.Proto,
			"service", "loms",
		)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
