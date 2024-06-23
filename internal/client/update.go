package client

import (
	"sync"

	model "github.com/gennadyterekhov/gophermart/internal/domain/models/order"
)

type Job struct {
	OrderNumber        string
	ResponseStatusCode int
	OrderStatus        string
	Accrual            *int64
	Error              error
}

var (
	jobsChannel chan *Job
	once        sync.Once
	RetryAfter  int64
	mu          sync.Mutex // maybe use atomics https://github.com/gennadyterekhov/gophermart/issues/24
)

func (ac *AccrualClient) initializeChannel() {
	jobsChannel = make(chan *Job)

	createWorkers()
	go func() {
		ac.workerPool()
	}()
}

func (ac *AccrualClient) LaunchAutoUpdate(order *model.Order) {
	once.Do(ac.initializeChannel)
	job := createJob(order)

	go func(job *Job) {
		jobsChannel <- job
	}(job)
}

func createJob(order *model.Order) *Job {
	job := &Job{
		OrderNumber: order.Number,
	}
	return job
}
