package delete

import (
	"context"

	"go.opentelemetry.io/otel"
)

type (
	repository interface {
		DeleteOne(ctx context.Context, userID, skuID int64)
	}

	Handler struct {
		repo repository
	}
)

func New(repo repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) DeleteItem(ctx context.Context, userID, skuID int64) {
	ctx, span := otel.Tracer("cart").Start(ctx, "service_delete_item")
	defer span.End()

	h.repo.DeleteOne(ctx, userID, skuID)
}
