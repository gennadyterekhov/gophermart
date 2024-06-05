package serializers

import (
	"encoding/json"

	"github.com/gennadyterekhov/gophermart/internal/domain/models/order"
)

func Orders(resDto *[]order.Order) ([]byte, error) {
	resDtoWithFloats := make([]order.OrderFloats, 0)

	for i := range *resDto {
		if (*resDto)[i].Accrual == nil {
			resDtoWithFloats = append(resDtoWithFloats, order.OrderFloats{
				Number:     (*resDto)[i].Number,
				Status:     (*resDto)[i].Status,
				UploadedAt: (*resDto)[i].UploadedAt,
			})
			continue
		}
		accrualFloat := float64(*(*resDto)[i].Accrual) / 100
		resDtoWithFloats = append(resDtoWithFloats, order.OrderFloats{
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
