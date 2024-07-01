package requests

type Withdrawals struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
