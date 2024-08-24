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
	"route256/cart/internal/repository/memorycartrepo"
	"route256/cart/pkg/prometheus"
)

type GetCartItemsResponse struct {
	Items      []domain.ListItem `json:"items"`
	TotalPrice int               `json:"total_price"`
}

type (
	getCartItemsCommand interface {
		GetItemsByUserID(ctx context.Context, userID int64) ([]domain.ListItem, error)
	}

	GetCartItemsHandler struct {
		name                string
		getCartItemsCommand getCartItemsCommand
	}

	getCartItemsRequest struct {
		// url params
		User int64 `validate:"nonzero"`
	}
)

func NewGetCartItemsHandler(command getCartItemsCommand, name string) *GetCartItemsHandler {
	return &GetCartItemsHandler{
		name:                name,
		getCartItemsCommand: command,
	}
}

func (h *GetCartItemsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	ctx, span := otel.Tracer("cart").Start(ctx, "handler_get_cart_items")
	defer span.End()

	defer func(createdAt time.Time) {
		prometheus.ObserveHttpRequestsDurationHistogram(createdAt, "get_cart_items")
	}(time.Now())

	prometheus.IncHttpRequestsTotalCounter("get_cart_items")

	var (
		request *getCartItemsRequest
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
			cartItems, err := h.getCartItemsCommand.GetItemsByUserID(ctx, request.User)
			if err != nil {
				if errors.Is(err, memorycartrepo.CartItemsNotFoundError{}) {
					GetErrorResponse(ctx, w, h.name, fmt.Errorf("command handler failed: %w", err), http.StatusNotFound)
					return
				}

				GetErrorResponse(ctx, w, h.name, fmt.Errorf("command handler failed: %w", err), http.StatusInternalServerError)

				return
			}

			response := GetCartItemsResponse{}
			response.Items = cartItems
			response.TotalPrice = getTotalPrice(cartItems)

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

func (_ *GetCartItemsHandler) getRequestData(r *http.Request) (request *getCartItemsRequest, err error) {
	request = &getCartItemsRequest{}

	if request.User, err = strconv.ParseInt(r.PathValue(definitions.ParamUserID), 10, 64); err != nil {
		return
	}

	return
}

func getTotalPrice(cartItems []domain.ListItem) int {
	var totalPrice uint32
	for _, item := range cartItems {
		totalPrice += item.Price * uint32(item.Count)
	}

	return int(totalPrice)
}
