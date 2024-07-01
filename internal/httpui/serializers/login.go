package serializers

import (
	"encoding/json"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
)

func Login(resDto *responses.Login) ([]byte, error) {
	serialized, err := json.Marshal(resDto)
	if err != nil {
		return []byte{}, err
	}

	return serialized, nil
}
