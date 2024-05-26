package serializers

import (
	"encoding/json"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"
)

func Withdrawals(resDto *[]models.Withdrawals) ([]byte, error) {
	serialized, err := json.Marshal(*resDto)
	if err != nil {
		return []byte{}, err
	}

	return serialized, nil
}
