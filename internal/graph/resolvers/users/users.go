package users

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"

	uuid "github.com/satori/go.uuid"

	"webinar/graphql/server/internal/graph/model"
	"webinar/graphql/server/internal/repositories/users"
	"webinar/graphql/server/internal/services/storage"
)

type UserResolvers interface {
	Login(ctx context.Context, email string, password string) (*model.Auth, error)
	CreateOneUser(ctx context.Context, input model.CreateUserInput) (*model.User, error)
	CreateOneCustomer(ctx context.Context, input model.CreateCustomerInput) (*model.Customer, error)
}

type resolvers struct {
	hashSalt       string
	sessionStorage storage.Session
	usersRepo      users.Repository
}

func New(usersRepo users.Repository, sessionStorage storage.Session) UserResolvers {
	return &resolvers{
		usersRepo:      usersRepo,
		sessionStorage: sessionStorage,
	}
}

func (r *resolvers) Login(ctx context.Context, email string, password string) (*model.Auth, error) {
	pwHash, err := passwordHash(password, r.hashSalt)
	if err != nil {
		return nil, err
	}

	session, err := r.usersRepo.FindUserSessionInfoByCredentials(ctx, email, pwHash)
	if err != nil {
		return nil, err
	}

	sessionToken := uuid.NewV4().String()
	if err := r.sessionStorage.Set(ctx, sessionToken, session); err != nil {
		return nil, fmt.Errorf("failed to set session key: %w", err)
	}

	return &model.Auth{Token: sessionToken}, nil
}

func (r *resolvers) CreateOneUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	pwHash, err := passwordHash(input.Password, r.hashSalt)
	if err != nil {
		return nil, err
	}
	input.Password = pwHash

	userId, err := r.usersRepo.InsertOrUpdateUser(ctx, nil, &input)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:    userId,
		Email: input.Email,
		Role:  input.Role,
	}, nil
}

func (r *resolvers) CreateOneCustomer(ctx context.Context, input model.CreateCustomerInput) (*model.Customer, error) {
	pwHash, err := passwordHash(input.Password, r.hashSalt)
	if err != nil {
		return nil, err
	}
	input.Password = pwHash

	customerId, userId, err := r.usersRepo.InsertOrUpdateCustomerUser(ctx, &input)
	if err != nil {
		return nil, err
	}

	return &model.Customer{
		ID:      customerId,
		UserID:  userId,
		Name:    input.Name,
		Address: input.Address,
	}, nil
}

func passwordHash(password, salt string) (string, error) {
	hash := md5.New()
	if _, err := io.WriteString(hash, password); err != nil {
		return "", fmt.Errorf("failed to add password to hash: %w", err)
	}
	if _, err := io.WriteString(hash, salt); err != nil {
		return "", fmt.Errorf("failed to add salt to hash: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
