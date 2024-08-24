package checkout

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"

	"route256/cart/internal/clients/loms"
	"route256/cart/internal/domain"
	"route256/cart/internal/repository/memorycartrepo"
)

type (
	lomsService interface {
		CreateOrder(ctx context.Context, userID int64, items []domain.Item) (int, error)
	}

	repository interface {
		GetAll(ctx context.Context, userID int64) ([]domain.Item, error)
		DeleteAll(_ context.Context, userID int64)
	}

	Handler struct {
		repo        repository
		lomsService lomsService
	}
)

func New(repo repository, lomsService lomsService) *Handler {
	return &Handler{
		repo:        repo,
		lomsService: lomsService,
	}
}

func (h *Handler) CartCheckout(ctx context.Context, userID int64) (*int, error) {
	ctx, span := otel.Tracer("cart").Start(ctx, "service_cart_checkout")
	defer span.End()

	cartItems, err := h.repo.GetAll(ctx, userID)

	if err != nil {
		if errors.Is(err, memorycartrepo.CartItemsNotFoundError{}) {
			return nil, memorycartrepo.CartItemsNotFoundError{}
		}

		return nil, fmt.Errorf("repository.GetCart: %w", err)
	}

	orderID, err := h.lomsService.CreateOrder(ctx, userID, cartItems)
	if err != nil {
		return nil, fmt.Errorf("%w %w", loms.ErrCreateOrder, err)
	}

	h.repo.DeleteAll(ctx, userID)

	return &orderID, nil
}
