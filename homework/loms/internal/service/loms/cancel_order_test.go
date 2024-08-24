package lomsusecase

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	orderStatus "route256/loms/internal/app/definitions"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/memory/orders"
	"route256/loms/internal/service/loms/mock"
)

func TestCancelOrderWithPrepare(t *testing.T) {
	ctx := context.Background()

	type (
		fields struct {
			ordersRepMock *mock.OrdersRepositoryMock
			stocksRepMock *mock.StocksRepositoryMock
		}

		data struct {
			name    string
			orderID int64
			prepare func(f *fields)
			wantErr error
		}
	)

	testData := []data{{
		name:    "Success",
		orderID: 123,
		prepare: func(f *fields) {
			orderItems := []domain.Item{{
				SKU:   872821,
				Count: 8,
			}}
			f.ordersRepMock.GetByIDMock.ExpectOrderIDParam2(123).Return(&domain.Order{
				Status: orderStatus.AwaitingPayment,
				UserID: 321,
				Items:  orderItems,
			}, nil)
			f.stocksRepMock.ReserveCancelMock.ExpectItemsParam2(orderItems).Return(nil)
			f.ordersRepMock.SetStatusMock.ExpectOrderIDParam2(123).ExpectStatusParam3(orderStatus.Cancelled).Return(nil)

			f.stocksRepMock.ReserveCancelMock.Times(1)
			f.ordersRepMock.SetStatusMock.Times(1)
		},
		wantErr: nil,
	}, {
		name:    "Order not found",
		orderID: 123,
		prepare: func(f *fields) {
			f.ordersRepMock.GetByIDMock.ExpectOrderIDParam2(123).Return(nil, orders.OrderNotFoundError{})
		},
		wantErr: orders.OrderNotFoundError{},
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
			err := handler.CancelOrder(ctx, tt.orderID)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
