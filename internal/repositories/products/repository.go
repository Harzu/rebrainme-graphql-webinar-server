package products

import (
	"context"

	"github.com/jmoiron/sqlx"

	"webinar/graphql/server/internal/entities"
	"webinar/graphql/server/internal/system/database/psql"
)

type Repository interface {
	InsertOrUpdateProduct(ctx context.Context, name string, price int64) (int64, error)
	FindProductsByIds(ctx context.Context, productIds []int64) ([]entities.Product, error)
	FindProductsByOrdersMap(ctx context.Context, orderIds []int64) (map[int64][]entities.Product, error)
}

type repositoryDB struct {
	client *sqlx.DB
}

func NewRepository(dbRepository psql.Repository) Repository {
	return &repositoryDB{client: dbRepository.GetConnection()}
}
