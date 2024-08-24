package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"gopkg.in/validator.v2"

	"route256/cart/internal/app/definitions"
	"route256/cart/pkg/prometheus"
)

type (
	deleteItemCommand interface {
		DeleteItem(ctx context.Context, userID, skuID int64)
	}

	DeleteItemHandler struct {
		name              string
		deleteItemCommand deleteItemCommand
	}
	deleteItemRequest struct {
		// url params
		SKU  int64 `validate:"nonzero"`
		User int64 `validate:"nonzero"`
	}
)

func NewDeleteItemHandler(command deleteItemCommand, name string) *DeleteItemHandler {
	return &DeleteItemHandler{
		name:              name,
		deleteItemCommand: command,
	}
}

func (h *DeleteItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	ctx, span := otel.Tracer("cart").Start(ctx, "handler_delete_cart_item")
	defer span.End()

	defer func(createdAt time.Time) {
		prometheus.ObserveHttpRequestsDurationHistogram(createdAt, "delete_cart_item")
	}(time.Now())

	prometheus.IncHttpRequestsTotalCounter("delete_cart_item")

	var (
		request *deleteItemRequest
		err     error
	)

	if request, err = h.getRequestData(r); err != nil {
		GetErrorResponse(ctx, w, h.name, err, http.StatusBadRequest)
		return
	}

	if err = validator.Validate(request); err != nil {
		GetErrorResponse(ctx, w, h.name, err, http.StatusBadRequest)
		return
	}

	done := make(chan struct{})

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer close(done)
		defer cancel()

		select {
		case <-ctx.Done():
			GetErrorResponse(ctx, w, h.name, fmt.Errorf("request context done: %w", ctx.Err()), http.StatusInternalServerError)
			return
		default:
			h.deleteItemCommand.DeleteItem(
				ctx,
				request.User,
				request.SKU,
			)

			GetNoContentResponse(ctx, w, h.name)
		}
	}()

	select {
	case <-done:
		return
	case <-r.Context().Done():
		GetErrorResponse(ctx, w, h.name, fmt.Errorf("request context done while waiting for goroutine: %w", r.Context().Err()), http.StatusInternalServerError)
		return
	}
}

func (_ *DeleteItemHandler) getRequestData(r *http.Request) (request *deleteItemRequest, err error) {
	request = &deleteItemRequest{}

	if request.User, err = strconv.ParseInt(r.PathValue(definitions.ParamUserID), 10, 64); err != nil {
		return
	}

	if request.SKU, err = strconv.ParseInt(r.PathValue(definitions.ParamSkuID), 10, 64); err != nil {
		return
	}

	return
}
