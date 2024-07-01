package controllers

import (
	"github.com/gennadyterekhov/gophermart/internal/domain/services"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/balance"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/login"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/orders"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/register"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/withdrawals"
)

type Controllers struct {
	Withdrawals withdrawals.Controller
	Orders      orders.Controller
	Balance     balance.Controller
	Register    register.Controller
	Login       login.Controller
}

func NewControllers(servs *services.Services) *Controllers {
	return &Controllers{
		Withdrawals: withdrawals.NewController(servs.Withdrawals),
		Orders:      orders.NewController(servs.Orders),
		Balance:     balance.NewController(servs.Balance),
		Register:    register.NewController(servs.Register),
		Login:       login.NewController(servs.Login),
	}
}
