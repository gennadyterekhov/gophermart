package models

import "time"

type Withdrawals struct {
	ID          int64
	UserID      int64
	OrderNumber int64
	TotalSum    int64
	ProcessedAt time.Time
}
