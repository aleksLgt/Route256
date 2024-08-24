package memorycartrepo

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"route256/cart/internal/domain"
)

func TestGetAllTable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type data struct {
		name       string
		userID     int64
		existItems itemsMap
		wantErr    error
	}

	testData := []data{{
		name:       "User does not have a cart",
		userID:     123,
		existItems: nil,
		wantErr:    CartItemsNotFoundError{},
	}, {
		name:       "User has an empty cart",
		userID:     938,
		existItems: itemsMap{},
		wantErr:    CartItemsNotFoundError{},
	}, {
		name:   "User has cart items",
		userID: 532,
		existItems: itemsMap{
			195: {
				SKU:   195,
				Count: 3,
			},
			865: {
				SKU:   865,
				Count: 8,
			},
		},
		wantErr: CartItemsNotFoundError{},
	}}

	for _, tt := range testData {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage := NewMemoryStorage()

			if tt.existItems != nil {
				storage.items[tt.userID] = tt.existItems
			}

			cartItems, err := storage.GetAll(ctx, tt.userID)
			if err != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.EqualError(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, len(tt.existItems), len(cartItems))
			}
		})
	}
}

func TestGetAllForDifferentUsersConcurrently(t *testing.T) {
	t.Parallel()

	var (
		countOfUsers = 100
		wg           = sync.WaitGroup{}
		ctx          = context.Background()
		storage      = NewMemoryStorage()
	)

	wg.Add(countOfUsers)

	for i := range countOfUsers {
		item := domain.Item{
			SKU:   int64(i),
			Count: uint16(i),
		}
		go func() {
			defer wg.Done()
			storage.Add(ctx, int64(i), item)
		}()
	}

	wg.Wait()

	require.Equal(t, len(storage.items), countOfUsers)

	wg.Add(countOfUsers)

	for i := range countOfUsers {
		go func() {
			defer wg.Done()

			_, err := storage.GetAll(ctx, int64(i))

			if err != nil {
				t.Errorf("Error when calling GetAll: %v", err)
			}
		}()
	}

	wg.Wait()

	require.Equal(t, len(storage.items), countOfUsers)
}

func TestGetAllForOneUserConcurrently(t *testing.T) {
	t.Parallel()

	var (
		countOfUserItems          = 1
		countOfConcurrentRequests = 100
		item                      = domain.Item{
			SKU:   int64(323),
			Count: uint16(6),
		}
		wg      = sync.WaitGroup{}
		ctx     = context.Background()
		storage = NewMemoryStorage()
		userID  = int64(123)
	)

	wg.Add(countOfUserItems)

	go func() {
		defer wg.Done()
		storage.Add(ctx, userID, item)
	}()

	wg.Wait()

	require.Equal(t, len(storage.items[userID]), countOfUserItems)

	wg.Add(countOfConcurrentRequests)

	for range countOfConcurrentRequests {
		go func() {
			defer wg.Done()

			_, err := storage.GetAll(ctx, userID)

			if err != nil {
				t.Errorf("Error when calling GetAll: %v", err)
			}
		}()
	}

	wg.Wait()

	require.Equal(t, len(storage.items[userID]), countOfUserItems)
}

func BenchmarkGetAll(b *testing.B) {
	ctx := context.Background()
	storage := NewMemoryStorage()

	var (
		userID     = int64(938)
		existItems = itemsMap{
			195: {
				SKU:   195,
				Count: 3,
			},
			865: {
				SKU:   865,
				Count: 8,
			},
			726: {
				SKU:   726,
				Count: 2,
			},
		}
	)

	storage.items[userID] = existItems

	b.ResetTimer()
	b.Run("GetAll", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := storage.GetAll(ctx, userID)
			if err != nil {
				b.Fatalf("GetAll failed: %v", err)
			}
		}
	})
	b.Run("GetAllOld", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := storage.GetAllOld(ctx, userID)
			if err != nil {
				b.Fatalf("GetAllOld failed: %v", err)
			}
		}
	})
}
