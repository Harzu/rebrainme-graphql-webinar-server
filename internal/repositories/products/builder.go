package products

import (
	"time"

	sq "github.com/Masterminds/squirrel"
)

const tableProducts = "products"

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
