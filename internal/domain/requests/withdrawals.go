package requests

type Withdrawals struct {
	Order string `json:"order"`
	Sum   int64  `json:"sum"`
}
