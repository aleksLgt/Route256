package add

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"route256/cart/internal/clients/loms"
	"route256/cart/internal/clients/product"
	"route256/cart/internal/domain"
	"route256/cart/internal/service/cart/item/add/mock"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestAddItemTableWithPrepare(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type (
		fields struct {
			productMock *mock.ProductServiceMock
			repMock     *mock.RepositoryMock
			lomsMock    *mock.LomsServiceMock
		}

		data struct {
			name            string
			userID          int64
			item            domain.Item
			prepare         func(f *fields)
			infoStocksCount int
			wantErr         error
		}
	)

	testData := []data{
		{
			name:   "product not found",
			userID: 123,
			item: domain.Item{
				SKU:   100,
				Count: 2,
			},
			prepare: func(f *fields) {
				f.productMock.GetProductInfoMock.ExpectSkuParam2(100).Return(nil, nil)
			},
			wantErr: ErrInvalidSKU,
		},
		{
			name:   "product service returned error",
			userID: 123,
			item: domain.Item{
				SKU:   111,
				Count: 5,
			},
			prepare: func(f *fields) {
				f.productMock.GetProductInfoMock.ExpectSkuParam2(111).Return(nil, fmt.Errorf("test error"))
			},
			wantErr: product.ErrGetProductInfo,
		},
		{
			name:   "loms service returned error",
			userID: 123,
			item: domain.Item{
				SKU:   100,
				Count: 2,
			},
			infoStocksCount: 4,
			prepare: func(f *fields) {
				f.productMock.GetProductInfoMock.ExpectSkuParam2(100).Return(&domain.Product{
					Name:  "Книга",
					Price: 300,
				}, nil)
				f.lomsMock.InfoStocksMock.ExpectSKUParam2(100).Return(0, fmt.Errorf("test error"))
			},
			wantErr: loms.ErrGetStockInfo,
		},
		{
			name:   "valid add item with enough stock",
			userID: 123,
			item: domain.Item{
				SKU:   100,
				Count: 2,
			},
			infoStocksCount: 4,
			prepare: func(f *fields) {
				f.productMock.GetProductInfoMock.ExpectSkuParam2(100).Return(&domain.Product{
					Name:  "Книга",
					Price: 300,
				}, nil)
				f.repMock.AddMock.ExpectUserIDParam2(123).ExpectItemParam3(domain.Item{
					SKU:   100,
					Count: 2,
				}).Return()
				f.lomsMock.InfoStocksMock.ExpectSKUParam2(100).Return(4, nil)
				f.repMock.AddMock.Times(1)
			},
			wantErr: nil,
		},
		{
			name:   "not enough stock",
			userID: 123,
			item: domain.Item{
				SKU:   100,
				Count: 2,
			},
			infoStocksCount: 4,
			prepare: func(f *fields) {
				f.productMock.GetProductInfoMock.ExpectSkuParam2(100).Return(&domain.Product{
					Name:  "Книга",
					Price: 300,
				}, nil)
				f.lomsMock.InfoStocksMock.ExpectSKUParam2(100).Return(1, nil)
			},
			wantErr: ErrInsufficientStocks,
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
				lomsMock:    mock.NewLomsServiceMock(ctrl),
			}

			addHandler := New(fieldsForTableTest.repMock, fieldsForTableTest.productMock, fieldsForTableTest.lomsMock)

			tt.prepare(&fieldsForTableTest)
			err := addHandler.AddItem(ctx, tt.userID, tt.item)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
