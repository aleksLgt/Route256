package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otelResource "go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"route256/cart/internal/app/closer"
	appHttp "route256/cart/internal/app/http"
	"route256/cart/internal/app/middleware"
	"route256/cart/internal/clients/loms"
	"route256/cart/internal/clients/product"
	"route256/cart/internal/domain"
	"route256/cart/internal/repository/memorycartrepo"
	cartCheckout "route256/cart/internal/service/cart/checkout"
	cartDelete "route256/cart/internal/service/cart/delete"
	cartItemAdd "route256/cart/internal/service/cart/item/add"
	cartItemDelete "route256/cart/internal/service/cart/item/delete"
	cartList "route256/cart/internal/service/cart/list"
	"route256/cart/pkg/logger"
)

type (
	mux interface {
		Handle(pattern string, handler http.Handler)
	}

	server interface {
		ListenAndServe() error
		Close() error
		Shutdown(ctx context.Context) error
	}

	cartStorage interface {
		Add(_ context.Context, userID int64, item domain.Item)
		DeleteOne(_ context.Context, userID, skuID int64)
		DeleteAll(_ context.Context, userID int64)
		GetAll(_ context.Context, userID int64) ([]domain.Item, error)
	}

	productClient interface {
		GetProductInfo(ctx context.Context, sku uint32) (*domain.Product, error)
	}

	lomsClient interface {
		CreateOrder(ctx context.Context, userID int64, items []domain.Item) (int, error)
		InfoStocks(ctx context.Context, SKU int64) (int, error)
	}

	App struct {
		ctx           context.Context
		config        *Config
		mux           mux
		server        server
		storage       cartStorage
		products      productClient
		lomsClient    lomsClient
		closer        *closer.Closer
		traceProvider *trace.TracerProvider
	}
)

func NewApp(ctx context.Context, config *Config) (*App, error) {
	mux := http.NewServeMux()

	_, err := logger.New()
	if err != nil {
		panic(err)
	}

	loggerCustom, err := logger.With("service", "cart")
	if err != nil {
		// panic, as an error can only occur if the logger is not initialized
		panic(err)
	}

	ctx = logger.ToContext(ctx, loggerCustom)

	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(config.jaegerAddr))
	if err != nil {
		logger.Panicw(ctx, "otel exporter error", "error", err)
	}

	resource, err := otelResource.Merge(
		otelResource.Default(),
		otelResource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("cart"),
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

	newProductsClient, err := product.New(config.productAddr, config.productToken)
	if err != nil {
		return nil, fmt.Errorf("the creation of a new product client failed: %w", err)
	}

	newLomsClient, err := loms.NewClient("user", config.lomsAddr)
	if err != nil {
		return nil, fmt.Errorf("the creation of a new loms client failed: %w", err)
	}

	return &App{
		ctx:    ctx,
		config: config,
		mux:    mux,
		server: &http.Server{
			Addr:              config.addr,
			Handler:           middleware.Logging(mux),
			ReadHeaderTimeout: 3 * time.Second,
		},
		storage:       memorycartrepo.NewMemoryStorage(),
		products:      newProductsClient,
		lomsClient:    newLomsClient,
		closer:        &closer.Closer{},
		traceProvider: traceProvider,
	}, nil
}

func (a *App) ListenAndServe() error {
	a.mux.Handle(a.config.path.cartItemAdd, appHttp.NewAddItemHandler(cartItemAdd.New(a.storage, a.products, a.lomsClient), a.config.path.cartItemAdd))
	a.mux.Handle(a.config.path.cartItemDelete, appHttp.NewDeleteItemHandler(cartItemDelete.New(a.storage), a.config.path.cartItemDelete))
	a.mux.Handle(a.config.path.cartDelete, appHttp.NewClearCartItemsHandler(cartDelete.New(a.storage), a.config.path.cartDelete))
	a.mux.Handle(a.config.path.cartList, appHttp.NewGetCartItemsHandler(cartList.New(a.storage, a.products), a.config.path.cartList))
	a.mux.Handle(a.config.path.cartCheckout, appHttp.NewCartCheckoutHandler(cartCheckout.New(a.storage, a.lomsClient), a.config.path.cartCheckout))
	a.mux.Handle(a.config.path.metrics, promhttp.Handler())
	a.mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	a.mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	a.mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	a.mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	a.mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	a.closer.Add(a.server.Shutdown)

	a.closer.Add(func(ctx context.Context) error {
		err := a.traceProvider.Shutdown(ctx)
		if err != nil {
			logger.Errorw(a.ctx, "traceProvider.Shutdown() failed", "error", err)
		}

		return nil
	})

	// При вызове a.closer.Close после остановки server и traceProvider ждем несколько секунд для завершения сторонних процессов
	a.closer.Add(func(ctx context.Context) error {
		time.Sleep(3 * time.Second)

		return nil
	})

	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Panicw(a.ctx, "listen and serve failed", "error", err)
		}
	}()

	logger.Infow(a.ctx, fmt.Sprintf("listening on %s", a.config.addr))
	<-a.ctx.Done()

	logger.Infow(a.ctx, "shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.closer.Close(shutdownCtx); err != nil {
		return fmt.Errorf("closer failed: %v", err)
	}

	return nil
}
