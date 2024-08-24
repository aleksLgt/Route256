package list

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"route256/cart/internal/clients/product"
	"route256/cart/internal/domain"
	"route256/cart/internal/repository/memorycartrepo"
	"route256/cart/internal/service/cart/list/mock"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestGetItemsByUserIDWithRepoErrorsWithPrepare(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type fields struct {
		repMock     *mock.RepositoryMock
		productMock *mock.ProductServiceMock
	}

	type data struct {
		name    string
		userID  int64
		prepare func(f *fields)
		wantErr error
	}

	testData := []data{{
		name:   "Repo returns CartItemsNotFoundError",
		userID: 123,
		prepare: func(f *fields) {
			f.repMock.GetAllMock.ExpectUserIDParam2(123).Return(nil, memorycartrepo.CartItemsNotFoundError{})
		},
		wantErr: memorycartrepo.CartItemsNotFoundError{},
	}, {
		name:   "Repo returns an error different from CartItemsNotFoundError",
		userID: 123,
		prepare: func(f *fields) {
			f.repMock.GetAllMock.ExpectUserIDParam2(123).Return(nil, fmt.Errorf("test error"))
		},
		wantErr: fmt.Errorf("repository.GetCart: test error"),
	}}

	for _, tt := range testData {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)

			fieldsForTableTest := fields{
				repMock:     mock.NewRepositoryMock(ctrl),
				productMock: mock.NewProductServiceMock(ctrl),
			}

			getHandler := New(fieldsForTableTest.repMock, fieldsForTableTest.productMock)

			tt.prepare(&fieldsForTableTest)

			_, err := getHandler.GetItemsByUserID(ctx, tt.userID)
			require.EqualError(t, err, tt.wantErr.Error())
		})
	}
}

func TestGetItemsByUserIDWithPrepare(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type fields struct {
		productMock *mock.ProductServiceMock
		repMock     *mock.RepositoryMock
	}

	type data struct {
		name    string
		userID  int64
		prepare func(f *fields)
		wantErr error
	}

	testData := []data{
		{
			name:   "Product service error",
			userID: 123,
			prepare: func(f *fields) {
				f.productMock.GetProductInfoMock.ExpectSkuParam2(uint32(234)).Return(nil, product.ErrGetProductInfo)
				f.repMock.GetAllMock.ExpectUserIDParam2(123).Return([]domain.Item{
					{
						SKU:   234,
						Count: 7,
					},
				}, nil)
			},
			wantErr: product.ErrGetProductInfo,
		},
		{
			name:   "Success",
			userID: 525,
			prepare: func(f *fields) {
				f.productMock.GetProductInfoMock.ExpectSkuParam2(uint32(983)).Return(&domain.Product{
					Name:  "Книга",
					Price: 400,
				}, nil)
				f.repMock.GetAllMock.ExpectUserIDParam2(525).Return([]domain.Item{
					{
						SKU:   983,
						Count: 2,
					},
				}, nil)
			},
			wantErr: nil,
		},
	}

	for _, tt := range testData {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)
			fieldsForTableTest := fields{
				productMock: mock.NewProductServiceMock(ctrl),
				repMock:     mock.NewRepositoryMock(ctrl),
			}

			getHandler := New(fieldsForTableTest.repMock, fieldsForTableTest.productMock)

			tt.prepare(&fieldsForTableTest)
			_, err := getHandler.GetItemsByUserID(ctx, tt.userID)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
