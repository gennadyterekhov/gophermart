package balance

import (
	"context"
	"fmt"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"

	"github.com/gennadyterekhov/gophermart/internal/repositories"
)

type Service struct {
	Repository repositories.Repository
}

func NewService(repo repositories.Repository) Service {
	return Service{
		Repository: repo,
	}
}

func (service *Service) GetBalanceResponse(ctx context.Context) (*responses.Balance, error) {
	userID, ok := ctx.Value(middleware.ContextUserIDKey).(int64)
	if !ok {
		return nil, fmt.Errorf("cannot get user_id from context")
	}

	resDto := &responses.Balance{}

	balance, err := service.GetBalance(ctx, userID)
	if err != nil {
		return nil, err
	}
	withdrawn, err := service.GetWithdrawn(ctx, userID)
	if err != nil {
		return nil, err
	}

	resDto.Current = balance
	resDto.Withdrawn = withdrawn

	return resDto, nil
}

func (service *Service) GetBalance(ctx context.Context, userID int64) (int64, error) {
	orders, err := service.Repository.GetAllOrdersForUser(ctx, userID)
	if err != nil {
		return 0, err
	}
	withdrawn, err := service.GetWithdrawn(ctx, userID)
	if err != nil {
		return 0, err
	}

	var sum int64
	for _, order := range orders {
		if order.Accrual != nil {
			sum += *order.Accrual
		}
	}
	sum -= withdrawn

	return sum, nil
}

func (service *Service) GetWithdrawn(ctx context.Context, userID int64) (int64, error) {
	wdrs, err := service.Repository.GetAllWithdrawalsForUser(ctx, userID)
	if err != nil {
		return 0, err
	}
	var sum int64
	for _, wdr := range wdrs {
		sum += wdr.TotalSum
	}

	return sum, nil
}
