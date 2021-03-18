package users

import (
	"context"

	"github.com/jmoiron/sqlx"

	"webinar/graphql/server/internal/entities"
	"webinar/graphql/server/internal/graph/model"
	"webinar/graphql/server/internal/system/database/psql"
)

type Repository interface {
	InsertOrUpdateUser(ctx context.Context, tx *sqlx.Tx, createUserInput *model.CreateUserInput) (int64, error)
	InsertOrUpdateCustomerUser(ctx context.Context, createCustomerInput *model.CreateCustomerInput) (int64, int64, error)
	FindUserSessionInfoByCredentials(ctx context.Context, email, passwordHash string) (entities.Session, error)
	FindCustomerByUserId(ctx context.Context, userId int64) (entities.Customer, error)
}

type repositoryDB struct {
	client *sqlx.DB
}

func NewRepository(dbRepository psql.Repository) Repository {
	return &repositoryDB{client: dbRepository.GetConnection()}
}
