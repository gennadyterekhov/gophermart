package orders

import (
	"context"
	"fmt"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"
	"github.com/gennadyterekhov/gophermart/internal/domain/requests"
	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
)

const (
	ErrorNoContent = "no content"
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

func Create(ctx context.Context, reqDto *requests.Withdrawals) (*responses.PostWithdrawals, error) {
	// TODO https://github.com/gennadyterekhov/gophermart/issues/5
	return nil, nil
}
