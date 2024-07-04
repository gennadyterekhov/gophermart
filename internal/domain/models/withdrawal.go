package models

import "time"

type Withdrawal struct {
	ID          int64
	UserID      int64
	OrderNumber string
	TotalSum    int64
	ProcessedAt time.Time
}
