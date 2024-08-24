package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"route256/cart/internal/app"
	"route256/cart/pkg/logger"
)

const (
	defaultAddr        = ":8082"
	defaultProductAddr = "http://route256.pavl.uk:8080"
	defaultLOMSAddr    = "loms:50051"
	defaultJaegerAddr  = "http://jaeger:4318"

	productToken = "testtoken"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	opts := initOpts()

	service, err := app.NewApp(ctx, app.NewConfig(&opts))

	if err != nil {
		logger.Panicw(ctx, "error when creating a new app", "error", err)
	}

	err = service.ListenAndServe()
	if err != nil {
		log.Printf("error starting server: %s\n", err)
		panic(err)
	}
}

func initOpts() app.Options {
	options := app.Options{}

	flag.StringVar(&options.Addr, "addr", defaultAddr, fmt.Sprintf("server address, default: %q", defaultAddr))
	flag.StringVar(&options.ProductAddr, "product_addr", defaultProductAddr, fmt.Sprintf("products-service address, default: %q", defaultProductAddr))
	flag.StringVar(&options.LOMSAddr, "loms_addr", defaultLOMSAddr, fmt.Sprintf("loms-service address, default: %q", defaultLOMSAddr))
	flag.StringVar(&options.JaegerAddr, "jaeger_addr", defaultJaegerAddr, fmt.Sprintf("jaeger address, default: %q", defaultJaegerAddr))
	flag.StringVar(&options.ProductToken, "product_token", productToken, "products-service token")
	flag.Parse()

	return options
}
