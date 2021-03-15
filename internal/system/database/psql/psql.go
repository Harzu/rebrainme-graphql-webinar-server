package psql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const dialect = "postgres"

type Repository interface {
	GetConnection() *sqlx.DB
}

type repository struct {
	client *sqlx.DB
}

func New(dsn string) (Repository, error) {
	db, err := sqlx.Open(dialect, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connect to %s db: %w", dialect, err)
	}

	repo := repository{client: db}
	if err = repo.client.Ping(); err != nil {
		return &repo, fmt.Errorf("failed to check %s db connection: %w", dialect, err)
	}

	return &repo, nil
}

func (c repository) GetConnection() *sqlx.DB {
	return c.client
}
