package withdrawals

import (
	"context"
	"fmt"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"

	"github.com/gennadyterekhov/gophermart/internal/domain/requests"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"

	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
)

const (
	ErrorNoContent         = "no content"
	ErrorInsufficientFunds = "insufficient funds"
)

func GetAll(ctx context.Context) (*[]models.Withdrawal, error) {
	userID, ok := ctx.Value(middleware.ContextUserIDKey).(int64)
	if !ok {
		return nil, fmt.Errorf("cannot get user_id from context")
	}

	withdrawals, err := repositories.GetAllWithdrawalsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(withdrawals) == 0 {
		return nil, fmt.Errorf(ErrorNoContent)
	}

	return &withdrawals, nil
}

func Create(ctx context.Context, reqDto *requests.Withdrawals) (*responses.PostWithdrawals, error) {
	userID, ok := ctx.Value(middleware.ContextUserIDKey).(int64)
	if !ok {
		return nil, fmt.Errorf("cannot get user_id from context")
	}

	currentBalance, err := getBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	if reqDto.Sum > currentBalance {
		return nil, fmt.Errorf(ErrorInsufficientFunds)
	}

	_, err = repositories.AddWithdrawal(ctx, userID, reqDto.Order, reqDto.Sum, time.Time{})
	if err != nil {
		return nil, err
	}

	return &responses.PostWithdrawals{}, nil
}

func getBalance(ctx context.Context, userID int64) (int64, error) {
	orders, err := repositories.GetAllOrdersForUser(ctx, userID)
	if err != nil {
		return 0, err
	}
	var sum int64
	for _, order := range orders {
		if order.Accrual != nil {
			sum += *order.Accrual
		}
	}

	return sum, nil
}
