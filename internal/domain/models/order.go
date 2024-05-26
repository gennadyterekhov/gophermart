package models

import "time"

type Order struct {
	Number     string
	UserID     int64
	Status     string
	Accrual    *int64
	UploadedAt time.Time
}
