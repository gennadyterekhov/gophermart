package balance

import (
	"context"

	"github.com/gennadyterekhov/gophermart/internal/repositories"
)

func GetBalance(ctx context.Context, userID int64) (int64, error) {
	orders, err := repositories.GetAllOrdersForUser(ctx, userID)
	if err != nil {
		return 0, err
	}
	withdrawn, err := GetWithdrawn(ctx, userID)
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

func GetWithdrawn(ctx context.Context, userID int64) (int64, error) {
	wdrs, err := repositories.GetAllWithdrawalsForUser(ctx, userID)
	if err != nil {
		return 0, err
	}
	var sum int64
	for _, wdr := range wdrs {
		sum += wdr.TotalSum
	}

	return sum, nil
}
