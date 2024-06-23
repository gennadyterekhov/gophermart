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

type Service struct {
	Repository    repositories.Repository
	AccrualClient client.AccrualClient
}

func NewService(repo repositories.Repository, cl client.AccrualClient) Service {
	return Service{
		Repository:    repo,
		AccrualClient: cl,
	}
}

func (service *Service) GetAll(ctx context.Context) (*[]order.Order, error) {
	userID, ok := ctx.Value(middleware.ContextUserIDKey).(int64)
	if !ok {
		return nil, fmt.Errorf("cannot get user_id from context")
	}

	orders, err := service.Repository.GetAllOrdersForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, fmt.Errorf(ErrorNoContent)
	}

	return &orders, nil
}

func (service *Service) Create(ctx context.Context, reqDto *requests.Orders) error {
	var err error
	userID, ok := ctx.Value(middleware.ContextUserIDKey).(int64)
	if !ok {
		return fmt.Errorf("cannot get user_id from context")
	}
	var orderObj *order.Order
	orderObj, err = service.Repository.GetOrderByIDAndUserID(ctx, reqDto.Number, userID)
	if err == nil && orderObj != nil {
		return fmt.Errorf(ErrorNumberAlreadyUploaded)
	}

	orderObj, err = service.Repository.GetOrderByID(ctx, reqDto.Number)
	if err == nil && orderObj != nil {
		return fmt.Errorf(ErrorNumberAlreadyUploadedByAnotherUser)
	}

	orderObj, err = service.Repository.AddOrder(
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

	_, err = service.AccrualClient.RegisterOrderInAccrual(reqDto.Number)
	if err != nil {
		return err
	}

	err = service.Repository.UpdateOrder(
		ctx,
		reqDto.Number,
		order.Processing,
		nil,
	)
	if err != nil {
		return err
	}

	service.AccrualClient.LaunchAutoUpdate(orderObj)

	return nil
}
