package serializers

import (
	"encoding/json"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"
)

func Withdrawals(resDto *[]models.Withdrawal) ([]byte, error) {
	resDtoWithFloats := make([]responses.WithdrawalExternal, 0)

	for i := range *resDto {
		resDtoWithFloats = append(
			resDtoWithFloats,
			responses.WithdrawalExternal{
				ID:          (*resDto)[i].ID,
				UserID:      (*resDto)[i].UserID,
				OrderNumber: (*resDto)[i].OrderNumber,
				TotalSum:    float64((*resDto)[i].TotalSum) / 100,
				ProcessedAt: (*resDto)[i].ProcessedAt,
			},
		)
	}
	serialized, err := json.Marshal(resDtoWithFloats)
	if err != nil {
		return []byte{}, err
	}

	return serialized, nil
}
