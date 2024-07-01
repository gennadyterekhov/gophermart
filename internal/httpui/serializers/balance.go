package serializers

import (
	"encoding/json"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
)

func Balance(resDto *responses.Balance) ([]byte, error) {
	resDtoWithFloats := &responses.BalanceExternal{
		Current:   float64(resDto.Current) / 100,
		Withdrawn: float64(resDto.Withdrawn) / 100,
	}

	serialized, err := json.Marshal(resDtoWithFloats)
	if err != nil {
		return []byte{}, err
	}

	return serialized, nil
}
