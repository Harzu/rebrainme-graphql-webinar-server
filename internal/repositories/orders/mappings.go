package orders

import "webinar/graphql/server/internal/entities"

func buildOrderEntity(dbModel ordersModel) entities.Order {
	return entities.Order{
		ID:         dbModel.ID,
		CustomerID: dbModel.CustomerID,
		TotalPrice: dbModel.TotalPrice,
		Status:     dbModel.Status,
	}
}
