package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otelResource "go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"route256/loms/internal/app/closer"
	"route256/loms/internal/app/loms"
	"route256/loms/internal/jobs"
	"route256/loms/internal/mw"
	"route256/loms/internal/repository/db/orders"
	"route256/loms/internal/repository/db/stocks"
	lomsUsecase "route256/loms/internal/service/loms"
	desc "route256/loms/pkg/api/loms/v1"
	"route256/loms/pkg/logger"
	"route256/loms/pkg/producer"
)

func headerMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case "x-auth":
		return key, true
	default:
		return key, false
	}
}

/*
Порты, которые слушает сервер.
*/
const (
	grpcPortEnv       = "GRPC_PORT"
	httpPortEnv       = "HTTP_PORT"
	dbConnReadStrEnv  = "DB_CONN_READ"
	dbConnWriteStrEnv = "DB_CONN_WRITE"
	jaegerHost        = "JAEGER_HOST"
)

//go:embed assets
var assets embed.FS

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	_, err := logger.New()
	if err != nil {
		panic(err)
	}

	loggerCustom, err := logger.With("service", "loms")
	if err != nil {
		// panic, as an error can only occur if the logger is not initialized
		panic(err)
	}

	ctx = logger.ToContext(ctx, loggerCustom)

	closerC := &closer.Closer{}

	traceProvider := initTracerProvider(ctx)

	closerC.Add(func(ctx context.Context) error {
		err = traceProvider.Shutdown(ctx)
		if err != nil {
			logger.Errorw(ctx, "traceProvider.Shutdown() failed", "error", err)
		}

		return nil
	})

	prod, err := producer.New()

	if err != nil {
		logger.Panicw(ctx, "failed to init producer", "error", err)
	}

	closerC.Add(func(ctx context.Context) error {
		err = prod.Close()
		if err != nil {
			logger.Errorw(ctx, "failed to close sync producer", "error", err)
		}

		return nil
	})

	initEnv(ctx)

	grpcPort, err := strconv.Atoi(os.Getenv(grpcPortEnv))
	if err != nil {
		logger.Panicw(ctx, "failed to get grpcPort", "error", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		logger.Panicw(ctx, "failed to listen", "error", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			mw.Panic,
			mw.Logger,
			mw.Auth,
			mw.Validate,
		),
	)

	// Рефлексия - это возможность диктовать клиентам свой контракт
	reflection.Register(grpcServer)

	ctx, cancel := context.WithTimeout(ctx, time.Minute)

	defer cancel()

	connWrite, err := pgx.Connect(ctx, os.Getenv(dbConnWriteStrEnv))

	if err != nil {
		logger.Panicw(ctx, "failed to connect to write database", "error", err)
	}

	closerC.Add(func(ctx context.Context) error {
		err = connWrite.Close(ctx)
		if err != nil {
			logger.Errorw(ctx, "failed to close write connection with database", "error", err)
		}

		return nil
	})

	connRead, err := pgx.Connect(ctx, os.Getenv(dbConnReadStrEnv))
	if err != nil {
		logger.Panicw(ctx, "failed to connect to read database", "error", err)
	}

	closerC.Add(func(ctx context.Context) error {
		err = connRead.Close(ctx)
		if err != nil {
			logger.Errorw(ctx, "failed to close write connection with database", "error", err)
		}

		return nil
	})

	useCase := lomsUsecase.NewService(orders.NewStorage(connRead, connWrite), stocks.NewStorage(connRead, connWrite))
	controller := loms.NewService(useCase)

	job := jobs.InitJob(connRead, connWrite)

	closerC.Add(func(ctx context.Context) error {
		job.Shutdown()

		return nil
	})

	go func() {
		job.Run()
	}()

	// Сгенерированный метод из прото
	desc.RegisterLOMSServer(grpcServer, controller)

	// grpcServer.Serve блокирующий, поэтому создаем в отдельной горутине
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			logger.Panicw(ctx, "failed to serve", "error", err)
		}
	}()

	gwmux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(headerMatcher))
	initAdditionalRoutes(ctx, gwmux)

	httpPort, err := strconv.Atoi(os.Getenv(httpPortEnv))
	if err != nil {
		logger.Panicw(ctx, "failed to get grpcPort", "error", err)
	}

	if err = desc.RegisterLOMSHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf(":%d", grpcPort), []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}); err != nil {
		logger.Panicw(ctx, "failed to register gateway", "error", err)
	}

	gwServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", httpPort),
		Handler:           mw.WithHTTPLoggingMiddleware(gwmux),
		ReadHeaderTimeout: 3 * time.Second,
	}

	// При вызове closerC.Close после остановки server ждем несколько секунд для завершения сторонних процессов
	closerC.Add(func(ctx context.Context) error {
		time.Sleep(3 * time.Second)

		return nil
	})

	go func() {
		err = gwServer.ListenAndServe()
		if err != nil {
			logger.Panicw(ctx, "failed to ListenAndServe", "error", err)
		}
	}()

	<-ctx.Done()

	logger.Infow(ctx, "shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closerC.Close(shutdownCtx); err != nil {
		logger.Errorw(ctx, "closerC Close failed", "error", err)
	}
}

func initAdditionalRoutes(ctx context.Context, gwmux *runtime.ServeMux) {
	sUI, _ := fs.Sub(assets, "assets")
	err := gwmux.HandlePath("GET", "/api/*", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		http.StripPrefix("/api/", http.FileServerFS(sUI)).ServeHTTP(w, r)
	})

	if err != nil {
		logger.Panicw(ctx, "error by handling swagger", "error", err)
	}

	err = gwmux.HandlePath("GET", "/metrics", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	if err != nil {
		logger.Panicw(ctx, "error by handling metrics", "error", err)
	}

	err = gwmux.HandlePath("GET", "/debug/pprof", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		pprof.Index(w, r)
	})

	if err != nil {
		logger.Panicw(ctx, "error by handling GET /debug/pprof", "error", err)
	}

	err = gwmux.HandlePath("GET", "/debug/pprof/cmdline", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		pprof.Cmdline(w, r)
	})

	if err != nil {
		logger.Panicw(ctx, "error by handling GET /debug/pprof/cmdline", "error", err)
	}

	err = gwmux.HandlePath("GET", "/debug/pprof/profile", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		pprof.Profile(w, r)
	})

	if err != nil {
		logger.Panicw(ctx, "error by handling GET /debug/pprof/profile", "error", err)
	}

	err = gwmux.HandlePath("GET", "/debug/pprof/symbol", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		pprof.Symbol(w, r)
	})

	if err != nil {
		logger.Panicw(ctx, "error by handling GET /debug/pprof/symbol", "error", err)
	}

	err = gwmux.HandlePath("GET", "/debug/pprof/trace", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		pprof.Trace(w, r)
	})

	if err != nil {
		logger.Panicw(ctx, "error by handling GET /debug/pprof/trace", "error", err)
	}
}

func initTracerProvider(ctx context.Context) *trace.TracerProvider {
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(os.Getenv(jaegerHost)))
	if err != nil {
		logger.Panicw(ctx, "otel exporter error", "error", err)
	}

	resource, err := otelResource.Merge(
		otelResource.Default(),
		otelResource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("loms"),
			semconv.DeploymentEnvironment("development"),
		),
	)
	if err != nil {
		logger.Panicw(ctx, "creating resource return error", "error", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource),
	)

	otel.SetTracerProvider(traceProvider)

	return traceProvider
}

func initEnv(ctx context.Context) {
	err := godotenv.Load()

	if err != nil {
		logger.Panicw(ctx, "Error loading .env file", "error", err)
	}
}
