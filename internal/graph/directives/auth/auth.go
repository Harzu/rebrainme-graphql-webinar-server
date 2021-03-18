package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/99designs/gqlgen/graphql"

	"webinar/graphql/server/internal/constants"
	"webinar/graphql/server/internal/graph/model"
	"webinar/graphql/server/internal/repositories/users"
	"webinar/graphql/server/internal/services/storage"
)

type Directive interface {
	Resolve(ctx context.Context, obj interface{}, next graphql.Resolver, roles []model.Role) (res interface{}, err error)
}

type authDirective struct {
	sessionStorage storage.Session
	usersRepo      users.Repository
}

func New(usersRepo users.Repository, sessionStorage storage.Session) Directive {
	return &authDirective{
		sessionStorage: sessionStorage,
		usersRepo:      usersRepo,
	}
}

func (d *authDirective) Resolve(
	ctx context.Context,
	obj interface{},
	next graphql.Resolver,
	roles []model.Role,
) (interface{}, error) {
	sessionKey, ok := ctx.Value(constants.SessionHeader).(string)
	if !ok {
		return nil, errors.New("failed to get session token from context")
	}

	session, err := d.sessionStorage.Get(ctx, sessionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from session storage: %w", err)
	}

	for _, role := range roles {
		if role == session.Role {
			return next(context.WithValue(ctx, constants.SessionContextKey, session))
		}
	}

	return nil, errors.New("unauthorized")
}
