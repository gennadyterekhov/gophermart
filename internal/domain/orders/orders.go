package orders

import (
	"context"
	"fmt"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/client"

	"github.com/gennadyterekhov/gophermart/internal/domain/models/order"

	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
)

const (
	ErrorNoContent                          = "no content"
	ErrorNumberAlreadyUploaded              = "ErrorNumberAlreadyUploaded"
	ErrorNumberAlreadyUploadedByAnotherUser = "ErrorNumberAlreadyUploadedByAnotherUser"
)

func GetAll(ctx context.Context) (*[]order.Order, error) {
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
	var err error
	userID, ok := ctx.Value(middleware.ContextUserIDKey).(int64)
	if !ok {
		return fmt.Errorf("cannot get user_id from context")
	}
	var orderObj *order.Order
	orderObj, err = repositories.GetOrderByIDAndUserID(ctx, reqDto.Number, userID)
	if err == nil && orderObj != nil {
		return fmt.Errorf(ErrorNumberAlreadyUploaded)
	}

	orderObj, err = repositories.GetOrderByID(ctx, reqDto.Number)
	if err == nil && orderObj != nil {
		return fmt.Errorf(ErrorNumberAlreadyUploadedByAnotherUser)
	}

	orderObj, err = repositories.AddOrder(
		ctx,
		reqDto.Number,
		userID,
		order.New,
		nil,
		time.Time{},
	)
	if err != nil {
		return err
	}

	_, err = client.RegisterOrderInAccrual(reqDto.Number)
	if err != nil {
		return err
	}

	err = repositories.UpdateOrder(
		ctx,
		reqDto.Number,
		order.Processing,
		nil,
	)
	if err != nil {
		return err
	}

	client.LaunchAutoUpdate(orderObj)

	return nil
}
