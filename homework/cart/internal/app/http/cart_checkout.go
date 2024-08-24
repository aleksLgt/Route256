package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"gopkg.in/validator.v2"

	"route256/cart/pkg/prometheus"
)

type (
	cartCheckoutCommand interface {
		CartCheckout(ctx context.Context, userID int64) (*int, error)
	}

	CartCheckoutHandler struct {
		name                string
		cartCheckoutCommand cartCheckoutCommand
	}

	cartCheckoutRequest struct {
		// url params
		User int64 `json:"user" validate:"nonzero"`
	}

	cartCheckoutResponse struct {
		OrderID int `json:"orderID"`
	}
)

func NewCartCheckoutHandler(command cartCheckoutCommand, name string) *CartCheckoutHandler {
	return &CartCheckoutHandler{
		name:                name,
		cartCheckoutCommand: command,
	}
}

func (h *CartCheckoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	ctx, span := otel.Tracer("cart").Start(ctx, "handler_cart_checkout")
	defer span.End()

	defer func(createdAt time.Time) {
		prometheus.ObserveHttpRequestsDurationHistogram(createdAt, "cart_checkout")
	}(time.Now())

	prometheus.IncHttpRequestsTotalCounter("cart_checkout")

	var (
		request *cartCheckoutRequest
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
			orderID, err := h.cartCheckoutCommand.CartCheckout(ctx, request.User)
			if err != nil {
				GetErrorResponse(ctx, w, h.name, err, http.StatusBadRequest)
				return
			}

			response := cartCheckoutResponse{OrderID: *orderID}

			buf, err := json.Marshal(&response)
			if err != nil {
				GetErrorResponse(ctx, w, h.name, fmt.Errorf("failed to encode response %w", err), http.StatusInternalServerError)
			}

			GetSuccessResponseWithBody(ctx, w, buf, h.name)
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

func (_ *CartCheckoutHandler) getRequestData(r *http.Request) (*cartCheckoutRequest, error) {
	request := &cartCheckoutRequest{}
	err := json.NewDecoder(r.Body).Decode(request)

	if err != nil {
		return nil, fmt.Errorf("failed to decode request data: %w", err)
	}

	return request, nil
}
