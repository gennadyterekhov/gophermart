package repositories

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/domain/models/order"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
	users                map[int64]*models.User
	lastUsedUserID       int64
	lastUsedWithdrawalID int64
	orders               map[string]*order.Order
	withdrawals          map[int64]*models.Withdrawal
}

func (repo *RepositoryMock) Clear() {
	repo.users = make(map[int64]*models.User)
	repo.orders = make(map[string]*order.Order)
	repo.withdrawals = make(map[int64]*models.Withdrawal)
}

func NewRepositoryMock() *RepositoryMock {
	return &RepositoryMock{
		users:       make(map[int64]*models.User),
		orders:      make(map[string]*order.Order),
		withdrawals: make(map[int64]*models.Withdrawal),
	}
}

func (repo *RepositoryMock) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	user, ok := repo.users[id]
	if !ok {
		return nil, nil
	}

	return user, nil
}

func (repo *RepositoryMock) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	for _, v := range repo.users {
		if v.Login == login {
			return v, nil
		}
	}

	return nil, nil
}

func (repo *RepositoryMock) AddUser(ctx context.Context, login string, password string) (*models.User, error) {
	alreadyExisting, err := repo.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if alreadyExisting != nil {
		return nil, fmt.Errorf("ERROR: duplicate key value violates unique constraint \"users_login_key\" (SQLSTATE 23505)")
	}

	repo.lastUsedUserID += 1
	newID := repo.lastUsedUserID
	user := &models.User{
		ID:       newID,
		Login:    login,
		Password: password,
	}

	repo.users[newID] = user

	return user, nil
}

func (repo *RepositoryMock) GetAllOrdersForUser(ctx context.Context, userID int64) ([]order.Order, error) {
	ords := make([]order.Order, 0)
	for _, v := range repo.orders {
		if v.UserID == userID {
			ords = append(ords, *v)
		}
	}

	sort.Slice(ords, func(i, j int) bool {
		return ords[i].UploadedAt.Before(ords[j].UploadedAt)
	})
	return ords, nil
}

func (repo *RepositoryMock) GetOrderByID(ctx context.Context, number string) (*order.Order, error) {
	ord, ok := repo.orders[number]
	if !ok {
		return nil, nil
	}

	return ord, nil
}

func (repo *RepositoryMock) GetOrderByIDAndUserID(ctx context.Context, number string, userID int64) (*order.Order, error) {
	ord, ok := repo.orders[number]
	if !ok {
		return nil, nil
	}
	if ord == nil {
		return nil, nil
	}
	if ord.UserID == userID {
		return ord, nil
	}

	return nil, nil
}

func (repo *RepositoryMock) AddOrder(
	ctx context.Context,
	number string,
	userID int64,
	status string,
	accrual *int64,
	uploadedAt time.Time,
) (*order.Order, error) {
	ord := &order.Order{
		Number:     number,
		UserID:     userID,
		Status:     status,
		Accrual:    accrual,
		UploadedAt: uploadedAt,
	}
	repo.orders[number] = ord

	return ord, nil
}

func (repo *RepositoryMock) UpdateOrder(
	ctx context.Context,
	number string,
	status string,
	accrual *int64,
) error {
	ord, ok := repo.orders[number]
	if !ok {
		return fmt.Errorf("could not find order to update")
	}

	ord.Status = status
	ord.Accrual = accrual

	return nil
}

func (repo *RepositoryMock) GetAllWithdrawalsForUser(ctx context.Context, userID int64) ([]models.Withdrawal, error) {
	wdrs := make([]models.Withdrawal, 0)
	for _, v := range repo.withdrawals {
		if v.UserID == userID {
			wdrs = append(wdrs, *v)
		}
	}
	sort.Slice(wdrs, func(i, j int) bool {
		return wdrs[i].ProcessedAt.Before(wdrs[j].ProcessedAt)
	})
	return wdrs, nil
}

func (repo *RepositoryMock) AddWithdrawal(
	ctx context.Context,
	userID int64,
	orderNumber string,
	totalSum int64,
	processedAt time.Time,
) (*models.Withdrawal, error) {
	repo.lastUsedWithdrawalID += 1
	newID := repo.lastUsedWithdrawalID
	wdr := &models.Withdrawal{
		ID:          newID,
		UserID:      userID,
		OrderNumber: orderNumber,
		TotalSum:    totalSum,
		ProcessedAt: processedAt,
	}
	repo.withdrawals[newID] = wdr

	return wdr, nil
}
