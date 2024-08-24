package lomsusecase

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"route256/loms/internal/repository/memory/stocks"
	"route256/loms/internal/service/loms/mock"
)

func TestInfoStocksWithPrepare(t *testing.T) {
	ctx := context.Background()

	type (
		fields struct {
			ordersRepMock *mock.OrdersRepositoryMock
			stocksRepMock *mock.StocksRepositoryMock
		}

		data struct {
			name    string
			sku     uint32
			prepare func(f *fields)
			count   *int64
			wantErr error
		}
	)

	testData := []data{{
		name: "Success",
		sku:  123,
		prepare: func(f *fields) {
			count := int64(4)
			f.stocksRepMock.GetBySKUMock.ExpectSkuParam2(123).Return(&count, nil)
		},
		count:   func() *int64 { v := int64(4); return &v }(),
		wantErr: nil,
	}, {
		name: "Stock not found",
		sku:  123,
		prepare: func(f *fields) {
			f.stocksRepMock.GetBySKUMock.ExpectSkuParam2(123).Return(nil, stocks.StockNotFoundError{})
		},
		count:   nil,
		wantErr: stocks.StockNotFoundError{},
	}, {
		name: "Stock has already been reserved",
		sku:  123,
		prepare: func(f *fields) {
			count := int64(0)
			f.stocksRepMock.GetBySKUMock.ExpectSkuParam2(123).Return(&count, nil)
		},
		count:   func() *int64 { v := int64(0); return &v }(),
		wantErr: nil,
	}}

	ctrl := minimock.NewController(t)
	fieldsForTableTest := fields{
		ordersRepMock: mock.NewOrdersRepositoryMock(ctrl),
		stocksRepMock: mock.NewStocksRepositoryMock(ctrl),
	}

	handler := NewService(fieldsForTableTest.ordersRepMock, fieldsForTableTest.stocksRepMock)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(&fieldsForTableTest)
			count, err := handler.InfoStocks(ctx, tt.sku)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.count, count)
		})
	}
}
