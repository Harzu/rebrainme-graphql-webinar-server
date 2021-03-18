package users

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"webinar/graphql/server/internal/entities"
	"webinar/graphql/server/internal/graph/model"
	"webinar/graphql/server/internal/repositories/common"
)

func (r *repositoryDB) FindUserSessionInfoByCredentials(_ context.Context, email, passwordHash string) (entities.Session, error) {
	query, args, err := prepareFindUserIdByCredentialsQuery(email, passwordHash)
	if err != nil {
		return entities.Session{}, fmt.Errorf("failed to prepare FindUserSessionInfoByCredentials query: %s", err)
	}

	var session sessionModel
	row := r.client.QueryRowx(query, args...)
	if err := row.StructScan(&session); err != nil {
		return entities.Session{}, fmt.Errorf("failed to scan row to session: %w", err)
	}

	return buildSessionEntity(session), nil
}

func (r *repositoryDB) InsertOrUpdateUser(_ context.Context, tx *sqlx.Tx, createUserInput *model.CreateUserInput) (int64, error) {
	query, args, err := prepareInsertOrUpdateUserQuery(createUserInput)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare InsertOrUpdateUser query: %s", err)
	}

	var client sqlx.Queryer = r.client
	if tx != nil {
		client = tx
	}

	var userId int64
	returningRow := client.QueryRowx(query, args...)
	if err := returningRow.Scan(&userId); err != nil {
		return 0, fmt.Errorf("failed to scan returning row to userId: %w", err)
	}

	return userId, nil
}

func (r *repositoryDB) InsertOrUpdateCustomerUser(ctx context.Context, createCustomerInput *model.CreateCustomerInput) (_ int64, _ int64, err error) {
	tx, err := r.client.Beginx()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to begin InsertOrUpdateCustomerUser tx: %w", err)
	}
	defer func() {
		commitErr := common.CommitTransaction(ctx, tx, err, "InsertOrUpdateCustomerUser")
		if commitErr != nil && err == nil {
			err = commitErr
		}
	}()

	userId, err := r.InsertOrUpdateUser(ctx, tx, &model.CreateUserInput{
		Email:    createCustomerInput.Email,
		Password: createCustomerInput.Password,
		Role:     model.RoleCustomer,
	})
	if err != nil {
		return 0, 0, err
	}

	query, args, err := prepareInsertOrUpdateCustomerQuery(createCustomerInput.Name, createCustomerInput.Address, userId)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to prepare InsertOrUpdateCustomerUser: %w", err)
	}

	var customerId int64
	returningRow := tx.QueryRowx(query, args...)
	if err := returningRow.Scan(&customerId); err != nil {
		return 0, 0, fmt.Errorf("failed to scan returning row to customerId: %w", err)
	}

	return customerId, userId, err
}

func (r *repositoryDB) FindCustomerByUserId(_ context.Context, userId int64) (entities.Customer, error) {
	query, args, err := prepareFindCustomerByUserIdQuery(userId)
	if err != nil {
		return entities.Customer{}, fmt.Errorf("failed to prepare FindCustomerByUserId query: %w", err)
	}

	var customer customerModel
	row := r.client.QueryRowx(query, args...)
	if err := row.StructScan(&customer); err != nil {
		return entities.Customer{}, fmt.Errorf("failed to scan row to customer: %w", err)
	}

	return buildCustomerEntity(customer), nil
}
