package orders

import (
	"context"

	"github.com/jmoiron/sqlx"

	"webinar/graphql/server/internal/entities"
	"webinar/graphql/server/internal/system/database/psql"
)

type Repository interface {
	InsertOrder(ctx context.Context, customerID int64, productIds []int64, totalPrice int64) (int64, error)
	FindOrdersByCustomerId(ctx context.Context, customerId int64) (entities.Orders, error)
}

type repositoryDB struct {
	client *sqlx.DB
}

func NewRepository(dbRepository psql.Repository) Repository {
	return &repositoryDB{client: dbRepository.GetConnection()}
}
