package repositories

import (
	"context"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/domain/models/order"
)

func (repo *Repository) GetAllOrdersForUser(ctx context.Context, userID int64) ([]order.Order, error) {
	const query = `SELECT 
    			       number, user_id, status, accrual, uploaded_at
				   FROM orders 
				   WHERE user_id = $1
				   ORDER BY uploaded_at`
	rows, err := repo.DB.Connection.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := make([]order.Order, 0)

	for rows.Next() {
		order := order.Order{}
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

func (repo *Repository) GetOrderByID(ctx context.Context, number string) (*order.Order, error) {
	const query = `SELECT number, user_id, status, accrual, uploaded_at FROM orders WHERE number = $1`
	row := repo.DB.Connection.QueryRowContext(ctx, query, number)
	if row.Err() != nil {
		return nil, row.Err()
	}

	order := order.Order{}
	err := row.Scan(&(order.Number), &(order.UserID), &(order.Status), &(order.Accrual), &(order.UploadedAt))
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (repo *Repository) GetOrderByIDAndUserID(ctx context.Context, number string, userID int64) (*order.Order, error) {
	const query = `SELECT number, user_id, status, accrual, uploaded_at FROM orders WHERE number = $1 and user_id = $2`
	row := repo.DB.Connection.QueryRowContext(ctx, query, number, userID)
	if row.Err() != nil {
		return nil, row.Err()
	}

	order := order.Order{}
	err := row.Scan(&(order.Number), &(order.UserID), &(order.Status), &(order.Accrual), &(order.UploadedAt))
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (repo *Repository) AddOrder(
	ctx context.Context,
	number string,
	userID int64,
	status string,
	accrual *int64,
	uploadedAt time.Time,
) (*order.Order, error) {
	const query = `INSERT INTO orders (number, user_id, status, accrual, uploaded_at)
			values ($1, $2, $3, $4, $5) RETURNING number;`

	_, err := repo.DB.Connection.ExecContext(ctx, query, number, userID, status, accrual, uploadedAt)
	if err != nil {
		return nil, err
	}

	order := order.Order{
		Number:     number,
		UserID:     userID,
		Status:     status,
		Accrual:    accrual,
		UploadedAt: uploadedAt,
	}

	return &order, nil
}

func (repo *Repository) UpdateOrder(
	ctx context.Context,
	number string,
	status string,
	accrual *int64,
) error {
	const query = `UPDATE orders SET status = $1, accrual = $2 WHERE number = $3;`

	_, err := repo.DB.Connection.ExecContext(ctx, query, status, accrual, number)
	if err != nil {
		return err
	}

	return nil
}
