package models

import "time"

type Withdrawal struct {
	ID          int64
	UserID      int64
	OrderNumber string
	TotalSum    int64
	ProcessedAt time.Time
}

type WithdrawalExternal struct {
	ID          int64
	UserID      int64
	OrderNumber string
	TotalSum    float64
	ProcessedAt time.Time
}
