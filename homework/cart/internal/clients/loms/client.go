package loms

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	header string
	conn   *grpc.ClientConn
}

func NewClient(header, addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("failed to create new gRPC loms client: %w", err)
	}

	return &Client{
		header: header,
		conn:   conn,
	}, nil
}
