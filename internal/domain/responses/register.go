package responses

type Register struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
}
