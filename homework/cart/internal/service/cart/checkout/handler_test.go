package checkout

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"route256/cart/internal/clients/loms"
	"route256/cart/internal/domain"
	"route256/cart/internal/repository/memorycartrepo"
	"route256/cart/internal/service/cart/checkout/mock"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestCheckoutCartWithRepoErrorsWithPrepare(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type (
		fields struct {
			repMock  *mock.RepositoryMock
			lomsMock *mock.LomsServiceMock
		}

		data struct {
			name    string
			userID  int64
			prepare func(f *fields)
			wantErr error
		}
	)

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
				repMock:  mock.NewRepositoryMock(ctrl),
				lomsMock: mock.NewLomsServiceMock(ctrl),
			}

			checkoutHandler := New(fieldsForTableTest.repMock, fieldsForTableTest.lomsMock)

			tt.prepare(&fieldsForTableTest)
			_, err := checkoutHandler.CartCheckout(ctx, tt.userID)
			require.EqualError(t, err, tt.wantErr.Error())
		})
	}
}

func TestCheckoutCartWithPrepare(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type (
		fields struct {
			repMock  *mock.RepositoryMock
			lomsMock *mock.LomsServiceMock
		}

		data struct {
			name    string
			userID  int64
			prepare func(f *fields)
			wantErr error
		}
	)

	testData := []data{{
		name:   "loms service returned error",
		userID: 123,
		prepare: func(f *fields) {
			cartItems := []domain.Item{
				{
					SKU:   983,
					Count: 2,
				},
			}
			f.repMock.GetAllMock.ExpectUserIDParam2(123).Return(cartItems, nil)
			f.lomsMock.CreateOrderMock.ExpectUserIDParam2(123).ExpectItemsParam3(cartItems).Return(0, fmt.Errorf("test error"))
		},
		wantErr: loms.ErrCreateOrder,
	}, {
		name:   "Success",
		userID: 123,
		prepare: func(f *fields) {
			cartItems := []domain.Item{
				{
					SKU:   983,
					Count: 2,
				},
			}
			f.repMock.GetAllMock.ExpectUserIDParam2(123).Return(cartItems, nil)
			f.lomsMock.CreateOrderMock.ExpectUserIDParam2(123).ExpectItemsParam3(cartItems).Return(2, nil)
			f.repMock.DeleteAllMock.ExpectUserIDParam2(123).Return()

			f.repMock.DeleteAllMock.Times(1)
		},
		wantErr: nil,
	}}

	for _, tt := range testData {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)
			fieldsForTableTest := fields{
				repMock:  mock.NewRepositoryMock(ctrl),
				lomsMock: mock.NewLomsServiceMock(ctrl),
			}

			checkoutHandler := New(fieldsForTableTest.repMock, fieldsForTableTest.lomsMock)

			tt.prepare(&fieldsForTableTest)
			_, err := checkoutHandler.CartCheckout(ctx, tt.userID)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
