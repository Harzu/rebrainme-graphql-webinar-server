package products

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const (
	tableProducts      = "products"
	tableOrderProducts = "order_products"
)

func prepareInsertOrUpdateProduct(name string, price int64) (string, []interface{}, error) {
	psqlSq := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	rawRequest := psqlSq.Insert(tableProducts).
		Columns("name", "price", "created_at").
		Values(name, price, time.Now()).
		Suffix(`
			ON CONFLICT (name) DO UPDATE SET
				price      = EXCLUDED.price,
				updated_at = now(),
				deleted_at = NULL
			RETURNING id
		`)

	return rawRequest.ToSql()
}

func prepareFindProductsByIdsQuery(productIds []int64) (string, []interface{}, error) {
	psqlSq := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	rawRequest := psqlSq.Select("id", "name", "price").
		From(tableProducts).
		Where(sq.Eq{
			"id":         productIds,
			"deleted_at": nil,
		})

	return rawRequest.ToSql()
}

func prepareFindProductsByOrdersMapQuery(orderIds []int64) (string, []interface{}, error) {
	psqlSq := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	rawRequest := psqlSq.Select("p.id", "p.name", "p.price", "op.order_id").
		From(fmt.Sprintf("%s p", tableProducts)).
		Join(fmt.Sprintf("%s op on p.id = op.product_id", tableOrderProducts)).
		Where(sq.Eq{
			"op.order_id":   orderIds,
			"op.deleted_at": nil,
			"p.deleted_at":  nil,
		})

	return rawRequest.ToSql()
}
