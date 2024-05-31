package orders

import (
	"context"
	"fmt"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
)

// 200 — номер заказа уже был загружен этим пользователем;
// 202 — новый номер заказа принят в обработку;
// 409 — номер заказа уже был загружен другим пользователем;
const (
	ErrorNoContent                          = "no content"
	ErrorNumberAlreadyUploaded              = "ErrorNumberAlreadyUploaded"
	ErrorNumberAlreadyUploadedByAnotherUser = "ErrorNumberAlreadyUploadedByAnotherUser"
)

func GetAll(ctx context.Context) (*[]models.Order, error) {
	userID, ok := ctx.Value(middleware.ContextUserIDKey).(int64)
	if !ok {
		return nil, fmt.Errorf("cannot get user_id from context")
	}

	orders, err := repositories.GetAllOrdersForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, fmt.Errorf(ErrorNoContent)
	}

	return &orders, nil
}

func Create(ctx context.Context, reqDto *requests.Orders) error {
	// TODO https://github.com/gennadyterekhov/gophermart/issues/10
	var err error
	userID, ok := ctx.Value(middleware.ContextUserIDKey).(int64)
	if !ok {
		return fmt.Errorf("cannot get user_id from context")
	}
	var order *models.Order
	order, err = repositories.GetOrderByIdAndUserId(ctx, reqDto.Number, userID)
	if err == nil && order != nil {
		return fmt.Errorf(ErrorNumberAlreadyUploaded)
	}

	order, err = repositories.GetOrderById(ctx, reqDto.Number)
	if err == nil && order != nil {
		return fmt.Errorf(ErrorNumberAlreadyUploadedByAnotherUser)
	}

	return nil
}
