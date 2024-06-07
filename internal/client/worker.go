package client

import (
	"context"
	"fmt"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/domain/models/order"
	"github.com/gennadyterekhov/gophermart/internal/logger"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
)

type Worker struct{}

const numberOfWorkers = 4

var workers []Worker

func createWorkers() {
	workers = make([]Worker, 0)
	workers = append(workers, *createWorker())
}

func createWorker() *Worker {
	return &Worker{}
}

func handleJob(job *Job) error {
	if job == nil {
		return fmt.Errorf("job is nil in worker")
	}
	logger.CustomLogger.Debugln("job for order", job.OrderNumber)

	var retryAfter int64

	mu.Lock()
	retryAfter = RetryAfter
	mu.Unlock()

	if retryAfter > 0 {
		logger.CustomLogger.Debugln("sleeping for", retryAfter)
		time.Sleep(time.Duration(retryAfter) * time.Second)
	}

	response, err := GetStatus(job.OrderNumber)
	if err != nil {
		return err
	}

	if response.CorrectResponse != nil {
		mu.Lock()
		RetryAfter = 0
		mu.Unlock()

		logger.CustomLogger.Debugln(
			"request to accrual with order "+job.OrderNumber+" was successful. new status",
			response.CorrectResponse.Status,
		)

		job.OrderStatus = response.CorrectResponse.Status
		if response.CorrectResponse.Accrual != nil {
			intAccrual := int64(100.0 * (*response.CorrectResponse.Accrual))
			job.Accrual = &intAccrual
			err = repositories.UpdateOrder(context.Background(), job.OrderNumber, job.OrderStatus, job.Accrual)
		} else {
			err = repositories.UpdateOrder(context.Background(), job.OrderNumber, job.OrderStatus, nil)
		}

		if err != nil {
			logger.CustomLogger.Errorln(err.Error())
			jobsChannel <- job
			return err
		}

		if job.OrderStatus == order.Processing {
			jobsChannel <- job
		}
	}

	if response.TooManyRequestsResponse != nil {
		mu.Lock()
		RetryAfter = response.TooManyRequestsResponse.RetryAfter
		mu.Unlock()
		jobsChannel <- job
	}

	return nil
}

func workerPool() {
	for w := 0; w < numberOfWorkers; w++ {
		go func() {
			for j := range jobsChannel {
				err := handleJob(j)
				if err != nil {
					logger.CustomLogger.Errorln("error in worker", err.Error())
					return
				}
			}
		}()
	}
}
