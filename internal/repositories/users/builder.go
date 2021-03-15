package users

import (
	"time"

	sq "github.com/Masterminds/squirrel"

	"webinar/graphql/server/internal/graph/model"
)

const (
	tableUsers     = "users"
	tableCustomers = "customers"
)

func prepareFindUserIdByCredentialsQuery(email string, passwordHash string) (string, []interface{}, error) {
	psqlSq := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	rawRequest := psqlSq.Select("id", "role").
		From(tableUsers).
		Where(sq.Eq{
			"email":         email,
			"password_hash": passwordHash,
			"deleted_at":    nil,
		})

	return rawRequest.ToSql()
}

func prepareInsertOrUpdateUserQuery(createUserInput *model.CreateUserInput) (string, []interface{}, error) {
	psqlSq := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	rawRequest := psqlSq.Insert(tableUsers).
		Columns("email", "password_hash", "role", "created_at").
		Values(createUserInput.Email, createUserInput.Password, createUserInput.Role, time.Now()).
		Suffix(`
			ON CONFLICT (email) DO UPDATE SET
				password_hash = EXCLUDED.password_hash,
				role          = EXCLUDED.role,
				updated_at    = now(),
				deleted_at    = NULL
			RETURNING id
		`)

	return rawRequest.ToSql()
}

func prepareInsertOrUpdateCustomerQuery(name string, address string, userId int64) (string, []interface{}, error) {
	psqlSq := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	rawRequest := psqlSq.Insert(tableCustomers).
		Columns("name", "address", "user_id", "created_at").
		Values(name, address, userId, time.Now()).
		Suffix(`
			ON CONFLICT (user_id) DO UPDATE SET
				name       = EXCLUDED.name,
				address    = EXCLUDED.address,
				updated_at = now(),
				deleted_at = NULL
			RETURNING id
		`)

	return rawRequest.ToSql()
}
