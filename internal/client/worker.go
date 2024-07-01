package client

import (
	"context"
	"fmt"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/domain/models/order"
	"github.com/gennadyterekhov/gophermart/internal/logger"
)

const numberOfWorkers = 4

func (ac *AccrualClient) handleJob(job *Job) error {
	if job == nil {
		return fmt.Errorf("job is nil in worker")
	}
	logger.CustomLogger.Debugln("job for order", job.OrderNumber)

	var retryAfter int64

	ac.mu.Lock()
	retryAfter = ac.RetryAfter
	ac.mu.Unlock()

	if retryAfter > 0 {
		logger.CustomLogger.Debugln("sleeping for", retryAfter)
		time.Sleep(time.Duration(retryAfter) * time.Second)
	}

	response, err := ac.GetStatus(job.OrderNumber)
	if err != nil {
		return err
	}

	if response.CorrectResponse != nil {
		ac.mu.Lock()
		ac.RetryAfter = 0
		ac.mu.Unlock()

		logger.CustomLogger.Debugln(
			"request to accrual with order "+job.OrderNumber+" was successful. new status",
			response.CorrectResponse.Status,
		)

		job.OrderStatus = response.CorrectResponse.Status
		if response.CorrectResponse.Accrual != nil {
			intAccrual := int64(100.0 * (*response.CorrectResponse.Accrual))
			job.Accrual = &intAccrual
			err = ac.Repository.UpdateOrder(context.Background(), job.OrderNumber, job.OrderStatus, job.Accrual)
		} else {
			err = ac.Repository.UpdateOrder(context.Background(), job.OrderNumber, job.OrderStatus, nil)
		}

		if err != nil {
			logger.CustomLogger.Errorln(err.Error())
			ac.JobsChannel <- job
			return err
		}

		if job.OrderStatus == order.Processing {
			ac.JobsChannel <- job
		}
	}

	if response.NoContentResponse != nil {
		ac.mu.Lock()
		ac.RetryAfter = 0
		ac.mu.Unlock()

		logger.CustomLogger.Debugln(
			"request to accrual with order "+job.OrderNumber+" was 'no content'. new status",
			response.NoContentResponse.Status,
		)

		job.OrderStatus = response.NoContentResponse.Status
		err = ac.Repository.UpdateOrder(context.Background(), job.OrderNumber, job.OrderStatus, nil)
		if err != nil {
			logger.CustomLogger.Errorln(err.Error())
			ac.JobsChannel <- job
			return err
		}

		if job.OrderStatus == order.Processing {
			ac.JobsChannel <- job
		}
	}

	if response.TooManyRequestsResponse != nil {
		logger.CustomLogger.Debugln(
			"request to accrual with order "+job.OrderNumber+" was 'too many requests'. RetryAfter:",
			response.TooManyRequestsResponse.RetryAfter,
		)
		ac.mu.Lock()
		ac.RetryAfter = response.TooManyRequestsResponse.RetryAfter
		ac.mu.Unlock()
		ac.JobsChannel <- job
	}

	return nil
}

func (ac *AccrualClient) workerPool() {
	for w := 0; w < numberOfWorkers; w++ {
		go func() {
			for j := range ac.JobsChannel {
				err := ac.handleJob(j)
				if err != nil {
					logger.CustomLogger.Errorln("error in worker", err.Error())
					return
				}
			}
		}()
	}
}
