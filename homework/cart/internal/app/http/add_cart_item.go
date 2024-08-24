package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"gopkg.in/validator.v2"

	"route256/cart/internal/app/definitions"
	"route256/cart/internal/domain"
	"route256/cart/internal/service/cart/item/add"
	"route256/cart/pkg/prometheus"
)

type (
	addItemCommand interface {
		AddItem(ctx context.Context, userID int64, item domain.Item) error
	}

	AddItemHandler struct {
		name           string
		addItemCommand addItemCommand
	}

	addItemRequest struct {
		// request body
		Count uint16 `son:"count" validate:"nonzero"`

		// url params
		SKU  int64 `validate:"nonzero"`
		User int64 `validate:"nonzero"`
	}
)

func NewAddItemHandler(command addItemCommand, name string) *AddItemHandler {
	return &AddItemHandler{
		name:           name,
		addItemCommand: command,
	}
}

func (h *AddItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	ctx, span := otel.Tracer("cart").Start(ctx, "handler_add_cart_item")
	defer span.End()

	defer func(createdAt time.Time) {
		prometheus.ObserveHttpRequestsDurationHistogram(createdAt, "add_cart_item")
	}(time.Now())

	prometheus.IncHttpRequestsTotalCounter("add_cart_item")

	var (
		request *addItemRequest
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
			err = h.addItemCommand.AddItem(
				ctx,
				request.User,
				domain.Item{
					SKU:   request.SKU,
					Count: request.Count,
				},
			)

			if err != nil {
				if errors.Is(err, add.ErrInvalidSKU) || errors.Is(err, add.ErrInsufficientStocks) {
					GetErrorResponse(ctx, w, h.name, fmt.Errorf("command handler failed: %w", err), http.StatusPreconditionFailed)
					return
				}

				GetErrorResponse(ctx, w, h.name, fmt.Errorf("command handler failed: %w", err), http.StatusInternalServerError)

				return
			}

			GetSuccessResponse(ctx, w, h.name)
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

func (_ *AddItemHandler) getRequestData(r *http.Request) (request *addItemRequest, err error) {
	request = &addItemRequest{}
	if err = json.NewDecoder(r.Body).Decode(request); err != nil {
		return
	}

	if request.User, err = strconv.ParseInt(r.PathValue(definitions.ParamUserID), 10, 64); err != nil {
		return
	}

	if request.SKU, err = strconv.ParseInt(r.PathValue(definitions.ParamSkuID), 10, 64); err != nil {
		return
	}

	return
}
