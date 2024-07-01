package client

import (
	"github.com/gennadyterekhov/gophermart/internal/domain/models"
)

type Job struct {
	OrderNumber        string
	ResponseStatusCode int
	OrderStatus        string
	Accrual            *int64
	Error              error
}

func (ac *AccrualClient) LaunchAutoUpdate(order *models.Order) {
	job := &Job{
		OrderNumber: order.Number,
	}

	go func(job *Job) {
		ac.JobsChannel <- job
	}(job)
}
