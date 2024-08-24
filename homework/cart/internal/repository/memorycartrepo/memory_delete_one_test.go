package memorycartrepo

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"route256/cart/internal/domain"
)

func TestDeleteOneTable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type data struct {
		name      string
		userID    int64
		SKU       int64
		existItem *domain.Item
	}

	testData := []data{{
		name:      "User does not have a cart",
		userID:    123,
		SKU:       195,
		existItem: nil,
	}, {
		name:   "User has items in a cart",
		userID: 938,
		SKU:    878,
		existItem: &domain.Item{
			SKU:   878,
			Count: 3,
		},
	}}

	for _, tt := range testData {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage := NewMemoryStorage()

			if tt.existItem != nil {
				storage.items[tt.userID] = itemsMap{tt.existItem.SKU: *tt.existItem}
			}

			storage.DeleteOne(ctx, tt.userID, tt.SKU)
			userItems := storage.items[tt.userID]

			_, ok := userItems[tt.SKU]
			if ok {
				t.Errorf("Item with SKU %d was not deleted from the storage", tt.SKU)
			}
		})
	}
}

func BenchmarkDeleteOne(b *testing.B) {
	ctx := context.Background()
	storage := NewMemoryStorage()

	var (
		userID    = int64(123)
		skuID     = int64(456)
		existItem = domain.Item{
			SKU:   456,
			Count: 3,
		}
	)

	storage.items[userID] = itemsMap{skuID: existItem}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		storage.DeleteOne(ctx, userID, skuID)
	}
}

func TestDeleteOneConcurrently(t *testing.T) {
	t.Parallel()

	var (
		countOfItems = 100
		ctx          = context.Background()
		storage      = NewMemoryStorage()
		userID       = int64(123)
		wg           = sync.WaitGroup{}
	)

	wg.Add(countOfItems)

	for i := range countOfItems {
		item := domain.Item{
			SKU:   int64(i),
			Count: uint16(i),
		}
		go func() {
			defer wg.Done()
			storage.Add(ctx, userID, item)
		}()
	}

	wg.Wait()

	require.Equal(t, len(storage.items[userID]), countOfItems)

	wg.Add(countOfItems)

	for sku := range countOfItems {
		go func() {
			defer wg.Done()
			storage.DeleteOne(ctx, userID, int64(sku))
		}()
	}

	wg.Wait()

	require.Equal(t, len(storage.items[userID]), 0)
}
