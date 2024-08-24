package memorycartrepo

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"route256/cart/internal/domain"
)

func TestDeleteAll(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	storage := NewMemoryStorage()

	var (
		userID     = int64(123)
		existItems = []domain.Item{
			{
				SKU:   195,
				Count: 3,
			}, {
				SKU:   865,
				Count: 8,
			},
		}
	)

	storage.items[userID] = itemsMap{existItems[0].SKU: existItems[0], existItems[1].SKU: existItems[1]}

	storage.DeleteAll(ctx, userID)
	_, ok := storage.items[userID]

	if ok {
		t.Errorf("Items for userID %d were not deleted from the storage", userID)
	}
}

func TestDeleteAllForDifferentUsersConcurrently(t *testing.T) {
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
			storage.DeleteAll(ctx, int64(i))
		}()
	}

	wg.Wait()

	require.Equal(t, len(storage.items), 0)
}

func TestDeleteAllForOneUserConcurrently(t *testing.T) {
	t.Parallel()

	var (
		countOfItems = 100
		wg           = sync.WaitGroup{}
		ctx          = context.Background()
		storage      = NewMemoryStorage()
		userID       = int64(123)
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

	for range countOfItems {
		go func() {
			defer wg.Done()
			storage.DeleteAll(ctx, userID)
		}()
	}

	wg.Wait()

	require.Equal(t, len(storage.items[userID]), 0)
}
