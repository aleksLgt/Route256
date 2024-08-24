package list

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	"go.opentelemetry.io/otel"
	"golang.org/x/time/rate"

	"route256/cart/internal/app/errgroup"
	"route256/cart/internal/clients/product"
	"route256/cart/internal/domain"
	"route256/cart/internal/repository/memorycartrepo"
)

type (
	productService interface {
		GetProductInfo(ctx context.Context, sku uint32) (*domain.Product, error)
	}
	repository interface {
		GetAll(ctx context.Context, userID int64) ([]domain.Item, error)
	}

	Handler struct {
		productService productService
		repo           repository
	}
)

func New(repo repository, productService productService) *Handler {
	return &Handler{
		repo:           repo,
		productService: productService,
	}
}

const requestsPerSecond = 10

func (h *Handler) GetItemsByUserID(ctx context.Context, userID int64) ([]domain.ListItem, error) {
	ctx, span := otel.Tracer("cart").Start(ctx, "service_get_items_by_user_id")
	defer span.End()

	// save to storage
	cartItems, err := h.repo.GetAll(ctx, userID)
	if err != nil {
		if errors.Is(err, memorycartrepo.CartItemsNotFoundError{}) {
			return nil, memorycartrepo.CartItemsNotFoundError{}
		}

		return nil, fmt.Errorf("repository.GetCart: %w", err)
	}

	sort.Slice(cartItems, func(i, j int) bool {
		return cartItems[i].SKU < cartItems[j].SKU
	})

	eg, ctx := errgroup.WithContext(ctx)

	limiter := rate.NewLimiter(requestsPerSecond, 1)

	listItemsChannel := make(chan domain.ListItem, len(cartItems))
	listItems := make([]domain.ListItem, len(cartItems))

	wg := sync.WaitGroup{}
	mx := sync.Mutex{}

	wg.Add(1)

	go func() {
		defer wg.Done()

		idx := 0

		for {
			select {
			case listItem, ok := <-listItemsChannel:
				if ok {
					mx.Lock()
					listItems[idx] = listItem
					mx.Unlock()

					idx++
				} else {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	for _, item := range cartItems {
		eg.Go(func() error {
			if err = limiter.Wait(ctx); err != nil {
				return fmt.Errorf("rate limiter error: %w", err)
			}

			productResponse, errGetProductInfo := h.productService.GetProductInfo(ctx, uint32(item.SKU))

			if errGetProductInfo != nil {
				return fmt.Errorf("%w %w", product.ErrGetProductInfo, err)
			}

			select {
			case listItemsChannel <- domain.ListItem{
				SKU:   item.SKU,
				Count: item.Count,
				Name:  productResponse.Name,
				Price: productResponse.Price,
			}:
			case <-ctx.Done():
				return ctx.Err()
			}

			return nil
		})
	}

	if err = eg.Wait(); err != nil {
		close(listItemsChannel)
		return nil, fmt.Errorf("%w", err)
	}

	wg.Wait()

	return listItems, nil
}

func (h *Handler) GetItemsByUserIDWithoutParallel(ctx context.Context, userID int64) ([]domain.ListItem, error) {
	// save to storage
	cartItems, err := h.repo.GetAll(ctx, userID)
	if err != nil {
		if errors.Is(err, memorycartrepo.CartItemsNotFoundError{}) {
			return nil, memorycartrepo.CartItemsNotFoundError{}
		}

		return nil, fmt.Errorf("repository.GetCart: %w", err)
	}

	sort.Slice(cartItems, func(i, j int) bool {
		return cartItems[i].SKU < cartItems[j].SKU
	})

	listItems := make([]domain.ListItem, len(cartItems))

	for i, item := range cartItems {
		productResponse, err := h.productService.GetProductInfo(ctx, uint32(item.SKU))
		if err != nil {
			return nil, fmt.Errorf("%w %w", product.ErrGetProductInfo, err)
		}

		listItems[i] = domain.ListItem{
			SKU:   item.SKU,
			Count: item.Count,
			Name:  productResponse.Name,
			Price: productResponse.Price,
		}
	}

	return listItems, nil
}
