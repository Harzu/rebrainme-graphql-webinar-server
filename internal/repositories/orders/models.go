package orders

import "webinar/graphql/server/internal/graph/model"

type ordersModel struct {
	ID         int64             `db:"id"`
	CustomerID int64             `db:"customer_id"`
	TotalPrice int64             `db:"total_price"`
	Status     model.OrderStatus `db:"status"`
}
