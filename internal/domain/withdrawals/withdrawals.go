package withdrawals

import (
	"context"
	"fmt"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/domain/balance"

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

	currentBalance, err := balance.GetBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	sumAsInt := int64(reqDto.Sum * 100)
	if sumAsInt > currentBalance {
		return nil, fmt.Errorf(ErrorInsufficientFunds)
	}

	_, err = repositories.AddWithdrawal(ctx, userID, reqDto.Order, sumAsInt, time.Time{})
	if err != nil {
		return nil, err
	}

	return &responses.PostWithdrawals{}, nil
}
