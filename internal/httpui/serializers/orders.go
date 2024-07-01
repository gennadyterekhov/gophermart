package serializers

import (
	"encoding/json"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"
)

func Orders(resDto *[]models.Order) ([]byte, error) {
	resDtoWithFloats := make([]models.OrderFloats, 0)

	for i := range *resDto {
		if (*resDto)[i].Accrual == nil {
			resDtoWithFloats = append(resDtoWithFloats, models.OrderFloats{
				Number:     (*resDto)[i].Number,
				Status:     (*resDto)[i].Status,
				UploadedAt: (*resDto)[i].UploadedAt,
			})
			continue
		}
		accrualFloat := float64(*(*resDto)[i].Accrual) / 100
		resDtoWithFloats = append(resDtoWithFloats, models.OrderFloats{
			Number:     (*resDto)[i].Number,
			Status:     (*resDto)[i].Status,
			UploadedAt: (*resDto)[i].UploadedAt,
			Accrual:    &accrualFloat,
		})
	}
	serialized, err := json.Marshal(resDtoWithFloats)
	if err != nil {
		return []byte{}, err
	}

	return serialized, nil
}
