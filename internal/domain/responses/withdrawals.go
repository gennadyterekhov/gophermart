package responses

import (
	"time"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"
)

type Withdrawals []models.Withdrawal

type WithdrawalExternal struct {
	ID          int64     `json:"-"`
	UserID      int64     `json:"-"`
	OrderNumber string    `json:"order"`
	TotalSum    float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type PostWithdrawals struct{}
