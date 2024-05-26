package repositories

import (
	"context"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"
	"github.com/gennadyterekhov/gophermart/internal/storage"
)

func GetAllWithdrawalsForUser(ctx context.Context, userID int64) ([]models.Withdrawal, error) {
	const query = `SELECT 
    			       id, 
    			       user_id, 
    			       order_number, 
    			       total_sum, 
    			       processed_at
				   FROM withdrawals 
				   WHERE user_id = $1
				   ORDER BY processed_at`
	rows, err := storage.DBClient.Connection.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	wdrs := make([]models.Withdrawal, 0)

	for rows.Next() {
		wdr := models.Withdrawal{}
		err = rows.Scan(&(wdr.ID), &(wdr.UserID), &(wdr.OrderNumber), &(wdr.TotalSum), &(wdr.ProcessedAt))
		if err != nil {
			return nil, err
		}
		wdrs = append(wdrs, wdr)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return wdrs, nil
}

func GetWithdrawalById(ctx context.Context, id int64) (*models.Withdrawal, error) {
	const query = `SELECT id, user_id, order_number, total_sum, processed_at FROM withdrawals WHERE id = $1`
	row := storage.DBClient.Connection.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	wdr := models.Withdrawal{}
	err := row.Scan(&(wdr.ID), &(wdr.UserID), &(wdr.OrderNumber), &(wdr.TotalSum), &(wdr.ProcessedAt))
	if err != nil {
		return nil, err
	}

	return &wdr, nil
}

func AddWithdrawal(
	ctx context.Context,
	userID int64,
	orderNumber string,
	totalSum int64,
	processedAt time.Time,
) (*models.Withdrawal, error) {
	const query = `INSERT INTO withdrawals (user_id, order_number, total_sum, processed_at)
			values ($1, $2, $3, $4) RETURNING id;`

	row := storage.DBClient.Connection.QueryRowContext(ctx, query, userID, orderNumber, totalSum, processedAt)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	wdr := models.Withdrawal{
		ID:          id,
		UserID:      userID,
		OrderNumber: orderNumber,
		TotalSum:    totalSum,
		ProcessedAt: processedAt,
	}

	return &wdr, nil
}
