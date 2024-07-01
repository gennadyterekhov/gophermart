package repositories

import (
	"context"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"
)

type RepositoryInterface interface {
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	AddUser(ctx context.Context, login string, password string) (*models.User, error)
	GetAllOrdersForUser(ctx context.Context, userID int64) ([]models.Order, error)
	GetOrderByID(ctx context.Context, number string) (*models.Order, error)
	GetOrderByIDAndUserID(ctx context.Context, number string, userID int64) (*models.Order, error)
	AddOrder(
		ctx context.Context,
		number string,
		userID int64,
		status string,
		accrual *int64,
		uploadedAt time.Time,
	) (*models.Order, error)
	UpdateOrder(
		ctx context.Context,
		number string,
		status string,
		accrual *int64,
	) error

	GetAllWithdrawalsForUser(ctx context.Context, userID int64) ([]models.Withdrawal, error)
	AddWithdrawal(
		ctx context.Context,
		userID int64,
		orderNumber string,
		totalSum int64,
		processedAt time.Time,
	) (*models.Withdrawal, error)
	Clear()
}
