package entities

import "webinar/graphql/server/internal/graph/model"

type Order struct {
	ID         int64
	CustomerID int64
	TotalPrice int64
	Status     model.OrderStatus
}

type Orders []Order

func (o Orders) ExtractIds() []int64 {
	result := make([]int64, 0, len(o))
	for _, order := range o {
		result = append(result, order.ID)
	}
	return result
}
