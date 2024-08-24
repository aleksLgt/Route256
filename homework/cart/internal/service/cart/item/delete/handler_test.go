package delete

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"go.uber.org/goleak"

	"route256/cart/internal/service/cart/item/delete/mock"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestDeleteItem(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type data struct {
		name   string
		userID int64
		skuID  int64
	}

	testData := data{
		name:   "Valid deleting of item",
		userID: 123,
		skuID:  796321,
	}

	ctrl := minimock.NewController(t)
	repMock := mock.NewRepositoryMock(ctrl)
	deleteHandler := New(repMock)

	repMock.DeleteOneMock.ExpectUserIDParam2(testData.userID).ExpectSkuIDParam3(testData.skuID).Return()
	deleteHandler.DeleteItem(ctx, testData.userID, testData.skuID)
	repMock.DeleteOneMock.Times(1)
}
