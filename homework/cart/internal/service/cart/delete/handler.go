package delete

import (
	"context"

	"go.opentelemetry.io/otel"
)

type (
	repository interface {
		DeleteAll(ctx context.Context, userID int64)
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

func (h *Handler) DeleteItemsByUserID(ctx context.Context, userID int64) {
	ctx, span := otel.Tracer("cart").Start(ctx, "service_delete_items_by_user_id")
	defer span.End()

	h.repo.DeleteAll(ctx, userID)
}
