package resolvers

import (
	"webinar/graphql/server/internal/graph/generated"
	"webinar/graphql/server/internal/graph/resolvers/customers"
	"webinar/graphql/server/internal/graph/resolvers/orders"
	"webinar/graphql/server/internal/graph/resolvers/products"
	"webinar/graphql/server/internal/graph/resolvers/users"
	"webinar/graphql/server/internal/repositories"
	"webinar/graphql/server/internal/services"
)

type container struct {
	users.UserResolvers
	orders.OrderResolvers
	products.ProductResolvers
	customers.CustomerResolvers
}

func New(services *services.Container, repositories *repositories.Container) *container {
	return &container{
		UserResolvers:     users.New(repositories.Users, services.SessionStorage),
		OrderResolvers:    orders.New(),
		ProductResolvers:  products.New(repositories.Products),
		CustomerResolvers: customers.New(),
	}
}

func (r *container) Mutation() generated.MutationResolver { return &mutation{r} }
func (r *container) Query() generated.QueryResolver       { return &query{r} }

type mutation struct{ *container }
type query struct{ *container }
