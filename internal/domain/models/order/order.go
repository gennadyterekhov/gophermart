package order

import "time"

// Order is used internally
type Order struct {
	Number     string
	UserID     int64
	Status     string
	Accrual    *int64
	UploadedAt time.Time
}

// OrderFloats is used to output
type OrderFloats struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *float64  `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func (ord *Order) isRegistered() bool {
	return ord.Status == Registered
}

func (ord *Order) isInvalid() bool {
	return ord.Status == Invalid
}

func (ord *Order) isProcessing() bool {
	return ord.Status == Processing
}

func (ord *Order) isProcessed() bool {
	return ord.Status == Processed
}
