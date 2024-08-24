package product

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"

	"route256/cart/internal/clients/product/middleware"
	"route256/cart/internal/domain"
	"route256/cart/pkg/prometheus"
)

type Client struct {
	token    string
	basePath string
}

type GetProductRequest struct {
	Token string `json:"token,omitempty"`
	SKU   uint32 `json:"sku,omitempty"`
}

type GetProductResponse struct {
	Name  string `json:"name,omitempty"`
	Price uint32 `json:"price,omitempty"`
}

type GetProductErrorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

const handlerName = "get_product"

var ErrGetProductInfo = errors.New("ProductService.GetProductInfo failed: ")

func New(basePath, token string) (*Client, error) {
	if token == "" {
		return nil, errors.New("product service has empty auth token")
	}

	return &Client{
		token:    token,
		basePath: basePath,
	}, nil
}

func (c Client) GetProductInfo(ctx context.Context, sku uint32) (*domain.Product, error) {
	ctx, span := otel.Tracer("cart").Start(ctx, "product_client_get_product_info")
	defer span.End()

	defer func(createdAt time.Time) {
		prometheus.ObserveExternalRequestsDurationHistogram(createdAt, "product", "get_product_info")
	}(time.Now())

	request := GetProductRequest{
		Token: c.token,
		SKU:   sku,
	}
	data, err := json.Marshal(request)

	if err != nil {
		return nil, fmt.Errorf("failed to encode request %w", err)
	}

	path, err := url.JoinPath(c.basePath, handlerName)
	if err != nil {
		return nil, fmt.Errorf("incorrect base basePath for %q: %w", handlerName, err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	client := &http.Client{
		Transport: &middleware.RetryTransport{MaxRetries: 3},
	}

	prometheus.IncExternalRequestsTotalCounter("product", "get_product_info")

	httpResponse, err := client.Do(httpRequest)
	prometheus.IncExternalResponseStatusTotalCounter("POST /get_product", strconv.Itoa(httpResponse.StatusCode))

	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	defer func() {
		_ = httpResponse.Body.Close()
	}()

	if httpResponse.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("product not found")
	}

	if httpResponse.StatusCode != http.StatusOK {
		response := &GetProductErrorResponse{}
		err = json.NewDecoder(httpResponse.Body).Decode(response)

		if err != nil {
			return nil, fmt.Errorf("failed to decode error response: %w", err)
		}

		return nil, fmt.Errorf("HTTP request responded with: %d , message: %s", httpResponse.StatusCode, response.Message)
	}

	response := &GetProductResponse{}
	err = json.NewDecoder(httpResponse.Body).Decode(response)

	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &domain.Product{
		Name:  response.Name,
		Price: response.Price,
	}, nil
}
