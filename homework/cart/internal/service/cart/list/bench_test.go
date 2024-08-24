package list

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"

	"route256/cart/internal/domain"
	"route256/cart/internal/service/cart/list/mock"
)

func BenchmarkGetItemsByUserID(b *testing.B) {
	ctx := context.Background()

	userID := int64(938)

	ctrl := minimock.NewController(b)
	repMock := mock.NewRepositoryMock(ctrl)
	productMock := mock.NewProductServiceMock(ctrl)
	listHandler := New(repMock, productMock)

	repMock.GetAllMock.Expect(ctx, userID).Return(getItems(), nil)

	productMock.GetProductInfoMock.Return(&domain.Product{
		Name:  "Книга",
		Price: 444,
	}, nil)

	b.ResetTimer()

	b.Run("GetItemsByUserID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := listHandler.GetItemsByUserID(ctx, userID)
			if err != nil {
				b.Fatalf("GetItemsByUserID failed: %v", err)
			}
		}
	})
}

func BenchmarkGetItemsByUserIDWithoutParallel(b *testing.B) {
	ctx := context.Background()

	userID := int64(938)

	ctrl := minimock.NewController(b)
	repMock := mock.NewRepositoryMock(ctrl)
	productMock := mock.NewProductServiceMock(ctrl)
	listHandler := New(repMock, productMock)

	repMock.GetAllMock.Expect(ctx, userID).Return(getItems(), nil)

	productMock.GetProductInfoMock.Return(&domain.Product{
		Name:  "Книга",
		Price: 444,
	}, nil)

	b.ResetTimer()

	b.Run("GetItemsByUserIDWithoutParallel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := listHandler.GetItemsByUserIDWithoutParallel(ctx, userID)
			if err != nil {
				b.Fatalf("GetItemsByUserIDWithoutParallel failed: %v", err)
			}
		}
	})
}

func getItems() []domain.Item {
	return []domain.Item{
		{
			SKU:   1076963,
			Count: 9,
		},
		{
			SKU:   1148162,
			Count: 5,
		},
		{
			SKU:   1625903,
			Count: 2,
		},
		{
			SKU:   2618151,
			Count: 4,
		},
		{
			SKU:   2956315,
			Count: 7,
		},
		{
			SKU:   2958025,
			Count: 6,
		},
		{
			SKU:   3596599,
			Count: 2,
		},
		{
			SKU:   3618852,
			Count: 1,
		},
		{
			SKU:   4288068,
			Count: 4,
		},
		{
			SKU:   4465995,
			Count: 3,
		},
	}
}
