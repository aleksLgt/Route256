package lomsusecase

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	orderStatus "route256/loms/internal/app/definitions"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/memory/stocks"
	"route256/loms/internal/service/loms/mock"
)

func TestCreateOrderWithPrepare(t *testing.T) {
	ctx := context.Background()

	type (
		fields struct {
			ordersRepMock *mock.OrdersRepositoryMock
			stocksRepMock *mock.StocksRepositoryMock
		}

		data struct {
			name       string
			userID     int64
			orderID    int64
			prepare    func(f *fields)
			orderItems []domain.Item
			wantErr    error
		}
	)

	testData := []data{{
		name:    "Success",
		userID:  123,
		orderID: 721,
		orderItems: []domain.Item{{
			SKU:   872821,
			Count: 8,
		}},
		prepare: func(f *fields) {
			orderItems := []domain.Item{{
				SKU:   872821,
				Count: 8,
			}}
			f.ordersRepMock.CreateMock.ExpectUserIDParam2(123).ExpectItemsParam3(orderItems).Return(721, nil)
			f.stocksRepMock.ReserveMock.ExpectItemsParam2(orderItems).Return(nil)
			f.ordersRepMock.SetStatusMock.ExpectOrderIDParam2(721).ExpectStatusParam3(orderStatus.AwaitingPayment).Return(nil)
		},
		wantErr: nil,
	}, {
		name:    "Stock not found",
		userID:  123,
		orderID: 721,
		orderItems: []domain.Item{{
			SKU:   872821,
			Count: 8,
		}},
		prepare: func(f *fields) {
			orderItems := []domain.Item{{
				SKU:   872821,
				Count: 8,
			}}
			f.ordersRepMock.CreateMock.ExpectUserIDParam2(123).ExpectItemsParam3(orderItems).Return(721, nil)
			f.stocksRepMock.ReserveMock.ExpectItemsParam2(orderItems).Return(stocks.StockNotFoundError{})
			f.ordersRepMock.SetStatusMock.ExpectOrderIDParam2(721).ExpectStatusParam3(orderStatus.Failed).Return(nil)
		},
		wantErr: stocks.StockNotFoundError{},
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
			_, err := handler.CreateOrder(ctx, tt.userID, tt.orderItems)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
