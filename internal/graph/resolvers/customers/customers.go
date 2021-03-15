package customers

import (
	"context"

	"webinar/graphql/server/internal/graph/model"
)

type CustomerResolvers interface {
	Me(ctx context.Context) (*model.Customer, error)
}

type resolvers struct{}

func New() CustomerResolvers {
	return &resolvers{}
}

func (r *resolvers) Me(ctx context.Context) (*model.Customer, error) {
	return &model.Customer{}, nil
}
