package repositories

import (
	"context"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"
	"github.com/gennadyterekhov/gophermart/internal/storage"
)

func GetAllOrdersForUser(ctx context.Context, userID int64) ([]models.Order, error) {
	const query = `SELECT 
    			       number, user_id, status, accrual, uploaded_at
				   FROM orders 
				   WHERE user_id = $1
				   ORDER BY uploaded_at`
	rows, err := storage.DBClient.Connection.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := make([]models.Order, 0)

	for rows.Next() {
		order := models.Order{}
		err = rows.Scan(&(order.Number), &(order.UserID), &(order.Status), &(order.Accrual), &(order.UploadedAt))
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func GetOrderById(ctx context.Context, number string) (*models.Order, error) {
	const query = `SELECT number, user_id, status, accrual, uploaded_at FROM orders WHERE number = $1`
	row := storage.DBClient.Connection.QueryRowContext(ctx, query, number)
	if row.Err() != nil {
		return nil, row.Err()
	}

	order := models.Order{}
	err := row.Scan(&(order.Number), &(order.UserID), &(order.Status), &(order.Accrual), &(order.UploadedAt))
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func GetOrderByIdAndUserId(ctx context.Context, number string, userID int64) (*models.Order, error) {
	const query = `SELECT number, user_id, status, accrual, uploaded_at FROM orders WHERE number = $1 and user_id = $2`
	row := storage.DBClient.Connection.QueryRowContext(ctx, query, number, userID)
	if row.Err() != nil {
		return nil, row.Err()
	}

	order := models.Order{}
	err := row.Scan(&(order.Number), &(order.UserID), &(order.Status), &(order.Accrual), &(order.UploadedAt))
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func AddOrder(
	ctx context.Context,
	number string,
	userID int64,
	status string,
	accrual *int64,
	uploadedAt time.Time,
) (*models.Order, error) {
	const query = `INSERT INTO orders (number, user_id, status, accrual, uploaded_at)
			values ($1, $2, $3, $4, $5) RETURNING number;`

	_, err := storage.DBClient.Connection.ExecContext(ctx, query, number, userID, status, accrual, uploadedAt)
	if err != nil {
		return nil, err
	}

	order := models.Order{
		Number:     number,
		UserID:     userID,
		Status:     status,
		Accrual:    accrual,
		UploadedAt: uploadedAt,
	}

	return &order, nil
}
