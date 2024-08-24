package delete

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"

	"route256/cart/internal/service/cart/delete/mock"
)

func TestDeleteItemsByUserID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type data struct {
		name   string
		userID int64
	}

	testData := data{
		name:   "Valid deleting of items",
		userID: 123,
	}

	ctrl := minimock.NewController(t)
	repMock := mock.NewRepositoryMock(ctrl)
	deleteHandler := New(repMock)

	repMock.DeleteAllMock.ExpectUserIDParam2(testData.userID).Return()
	deleteHandler.DeleteItemsByUserID(ctx, testData.userID)
	repMock.DeleteAllMock.Times(1)
}
