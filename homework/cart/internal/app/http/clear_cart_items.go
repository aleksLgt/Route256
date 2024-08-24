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
	clearCartItemsCommand interface {
		DeleteItemsByUserID(ctx context.Context, userID int64)
	}

	ClearCartItemsHandler struct {
		name                  string
		clearCartItemsCommand clearCartItemsCommand
	}
	clearCartItemsRequest struct {
		// url params
		User int64 `validate:"nonzero"`
	}
)

func NewClearCartItemsHandler(command clearCartItemsCommand, name string) *ClearCartItemsHandler {
	return &ClearCartItemsHandler{
		name:                  name,
		clearCartItemsCommand: command,
	}
}

func (h *ClearCartItemsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	ctx, span := otel.Tracer("cart").Start(ctx, "handler_clear_cart_item")
	defer span.End()

	defer func(createdAt time.Time) {
		prometheus.ObserveHttpRequestsDurationHistogram(createdAt, "clear_cart_items")
	}(time.Now())

	prometheus.IncHttpRequestsTotalCounter("clear_cart_items")

	var (
		request *clearCartItemsRequest
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
			h.clearCartItemsCommand.DeleteItemsByUserID(
				ctx,
				request.User,
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

func (_ *ClearCartItemsHandler) getRequestData(r *http.Request) (request *clearCartItemsRequest, err error) {
	request = &clearCartItemsRequest{}

	if request.User, err = strconv.ParseInt(r.PathValue(definitions.ParamUserID), 10, 64); err != nil {
		return
	}

	return
}
