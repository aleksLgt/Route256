package add

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"

	"route256/cart/internal/clients/loms"
	"route256/cart/internal/clients/product"
	"route256/cart/internal/domain"
)

type (
	productService interface {
		GetProductInfo(ctx context.Context, sku uint32) (*domain.Product, error)
	}

	repository interface {
		Add(ctx context.Context, userID int64, item domain.Item)
	}

	lomsService interface {
		InfoStocks(ctx context.Context, SKU int64) (int, error)
	}

	Handler struct {
		productService productService
		repo           repository
		lomsService    lomsService
	}
)

var (
	ErrInvalidSKU         = errors.New("invalid sku")
	ErrInsufficientStocks = errors.New("insufficient stocks")
)

func New(repo repository, productService productService, lomsService lomsService) *Handler {
	return &Handler{
		repo:           repo,
		productService: productService,
		lomsService:    lomsService,
	}
}

func (h *Handler) AddItem(ctx context.Context, userID int64, item domain.Item) error {
	ctx, span := otel.Tracer("cart").Start(ctx, "service_add_item")
	defer span.End()

	products, err := h.productService.GetProductInfo(ctx, uint32(item.SKU))
	if err != nil {
		return fmt.Errorf("%w %w", product.ErrGetProductInfo, err)
	}

	if products == nil {
		return fmt.Errorf("ProductService.GetProductInfo return no product with given SKU=%d: %w", item.SKU, ErrInvalidSKU)
	}

	count, err := h.lomsService.InfoStocks(ctx, item.SKU)
	if err != nil {
		return fmt.Errorf("%w %w", loms.ErrGetStockInfo, err)
	}

	if count < int(item.Count) {
		return fmt.Errorf("%w", ErrInsufficientStocks)
	}

	h.repo.Add(ctx, userID, item)

	return nil
}
