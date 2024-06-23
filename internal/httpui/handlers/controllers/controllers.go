package controllers

import (
	"github.com/gennadyterekhov/gophermart/internal/client"
	"github.com/gennadyterekhov/gophermart/internal/config"
	loginDomain "github.com/gennadyterekhov/gophermart/internal/domain/auth"
	registerDomain "github.com/gennadyterekhov/gophermart/internal/domain/auth/register"
	balanceDomain "github.com/gennadyterekhov/gophermart/internal/domain/balance"
	ordersDomain "github.com/gennadyterekhov/gophermart/internal/domain/orders"
	withdrawalsDomain "github.com/gennadyterekhov/gophermart/internal/domain/withdrawals"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/balance"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/login"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/orders"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/register"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/withdrawals"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
)

type Controllers struct {
	Withdrawals withdrawals.Controller
	Orders      orders.Controller
	Balance     balance.Controller
	Register    register.Controller
	Login       login.Controller
}

func NewControllers(conf *config.Config, repo repositories.RepositoryInterface) *Controllers {
	balanceService := balanceDomain.NewService(repo)
	return &Controllers{
		Withdrawals: withdrawals.NewController(withdrawalsDomain.NewService(repo, balanceService)),
		Orders:      orders.NewController(ordersDomain.NewService(repo, client.NewClient(conf.AccrualURL, repo))),
		Balance:     balance.NewController(balanceService),
		Register:    register.NewController(registerDomain.NewService(repo)),
		Login:       login.NewController(loginDomain.NewService(repo)),
	}
}
