package orders

import (
	"time"

	sq "github.com/Masterminds/squirrel"

	"webinar/graphql/server/internal/graph/model"
)

const (
	tableOrders        = "orders"
	tableOrderProducts = "order_products"
)

func prepareInsertOrderQuery(customerId int64, totalPrice int64) (string, []interface{}, error) {
	psqlSq := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	rawRequest := psqlSq.Insert(tableOrders).
		Columns("customer_id", "total_price", "status", "created_at").
		Values(customerId, totalPrice, model.OrderStatusCreated, time.Now()).
		Suffix("RETURNING id")

	return rawRequest.ToSql()
}

func prepareInsertOrderProductsQuery(orderId int64, productIds []int64) (string, []interface{}, error) {
	psqlSq := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	rawRequest := psqlSq.Insert(tableOrderProducts).
		Columns("order_id", "product_id", "created_at")

	for _, productId := range productIds {
		rawRequest = rawRequest.Values(orderId, productId, time.Now())
	}

	return rawRequest.ToSql()
}
