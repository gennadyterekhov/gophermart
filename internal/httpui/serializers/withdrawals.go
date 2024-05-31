package serializers

import (
	"encoding/json"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"
)

func Withdrawals(resDto *[]models.Withdrawal) ([]byte, error) {
	//resDtoWithFloats := make([]models.WithdrawalFloats, len(*resDto))
	//
	//for i := range *resDto {
	//	resDtoWithFloats = append(resDtoWithFloats, models.WithdrawalFloats{
	//		ID:          (*resDto)[i].ID,
	//		UserID:      (*resDto)[i].UserID,
	//		OrderNumber: (*resDto)[i].OrderNumber,
	//		TotalSum:    float64((*resDto)[i].TotalSum) / 100,
	//		ProcessedAt: (*resDto)[i].ProcessedAt,
	//	})
	//}
	//serialized, err := json.Marshal(resDtoWithFloats)
	// TODO serialize with floats https://github.com/gennadyterekhov/gophermart/issues/40
	serialized, err := json.Marshal(*resDto)
	if err != nil {
		return []byte{}, err
	}

	return serialized, nil
}
