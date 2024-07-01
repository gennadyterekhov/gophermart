package services

import (
	"github.com/gennadyterekhov/gophermart/internal/client"
	"github.com/gennadyterekhov/gophermart/internal/domain/auth"
	"github.com/gennadyterekhov/gophermart/internal/domain/auth/register"
	"github.com/gennadyterekhov/gophermart/internal/domain/balance"
	"github.com/gennadyterekhov/gophermart/internal/domain/orders"
	"github.com/gennadyterekhov/gophermart/internal/domain/withdrawals"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
)

type Services struct {
	Withdrawals withdrawals.Service
	Orders      orders.Service
	Balance     balance.Service
	Register    register.Service
	Login       auth.Service
}

func New(repo repositories.RepositoryInterface, accrualClient *client.AccrualClient) *Services {
	balanceService := balance.NewService(repo)
	return &Services{
		Withdrawals: withdrawals.NewService(repo, balanceService),
		Orders:      orders.NewService(repo, accrualClient),
		Balance:     balanceService,
		Register:    register.NewService(repo),
		Login:       auth.NewService(repo),
	}
}
