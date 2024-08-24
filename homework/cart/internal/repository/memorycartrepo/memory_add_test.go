package memorycartrepo

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"route256/cart/internal/domain"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestAddTable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type data struct {
		name        string
		userID      int64
		item        *domain.Item
		existItem   *domain.Item
		expectCount uint16
	}

	testData := []data{{
		name:      "User does not have an item with the same SKU",
		userID:    123,
		existItem: nil,
		item: &domain.Item{
			SKU:   111,
			Count: 5,
		},
		expectCount: 5,
	}, {
		name:   "User has an item with the same SKU",
		userID: 938,
		existItem: &domain.Item{
			SKU:   195,
			Count: 4,
		},
		item: &domain.Item{
			SKU:   195,
			Count: 5,
		},
		expectCount: 9,
	}}

	for _, tt := range testData {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage := NewMemoryStorage()

			if tt.existItem != nil {
				storage.items[tt.userID] = itemsMap{}
				storage.items[tt.userID][tt.existItem.SKU] = *tt.existItem
			}

			storage.Add(ctx, tt.userID, *tt.item)
			userItems := storage.items[tt.userID]

			addedItem, ok := userItems[tt.item.SKU]
			if !ok {
				t.Errorf("Item with SKU %d was not added to the storage", tt.item.SKU)
			}

			if addedItem.Count != tt.expectCount {
				t.Errorf("Incorrect item count. Expected: %d, Got: %d", tt.expectCount, addedItem.Count)
			}
		})
	}
}

func TestAddConcurrently(t *testing.T) {
	t.Parallel()

	var (
		countOfProducts = 100
		wg              = sync.WaitGroup{}
		ctx             = context.Background()
		userID          = int64(123)
		storage         = NewMemoryStorage()
	)

	wg.Add(countOfProducts)

	for i := range countOfProducts {
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

	userItems, err := storage.GetAll(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, countOfProducts, len(userItems))
}

func BenchmarkAdd(b *testing.B) {
	ctx := context.Background()
	storage := NewMemoryStorage()

	var (
		userID    = int64(938)
		existItem = domain.Item{
			SKU:   195,
			Count: 4,
		}
		item = domain.Item{
			SKU:   195,
			Count: 5,
		}
	)

	storage.items[userID] = itemsMap{existItem.SKU: existItem}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		storage.Add(ctx, userID, item)
	}
}
