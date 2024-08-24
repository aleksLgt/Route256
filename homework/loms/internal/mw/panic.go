package mw

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"route256/loms/pkg/logger"
)

func Panic(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if e := recover(); e != nil {
			logger.Errorw(ctx, "panic", "error", e)
			err = status.Errorf(codes.Internal, "panic: %v", e)
		}
	}()

	resp, err = handler(ctx, req)

	return resp, err
}
