package withdrawals

import (
	"context"
	"fmt"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"

	"github.com/gennadyterekhov/gophermart/internal/httpui/middleware"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
)

const ErrorNoContent = "no content"

func GetAll(ctx context.Context) (*[]models.Withdrawals, error) {
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
