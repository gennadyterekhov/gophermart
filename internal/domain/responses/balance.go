package responses

type Balance struct {
	Current   int64 `json:"current"`
	Withdrawn int64 `json:"withdrawn"`
}

type BalanceExternal struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
