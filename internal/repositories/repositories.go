package repositories

import (
	"webinar/graphql/server/internal/repositories/orders"
	"webinar/graphql/server/internal/repositories/products"
	"webinar/graphql/server/internal/repositories/users"
	"webinar/graphql/server/internal/system/database/psql"
)

// DI for repositories
type Container struct {
	Users    users.Repository
	Orders   orders.Repository
	Products products.Repository
}

func New(dbRepository psql.Repository) *Container {
	return &Container{
		Users:    users.NewRepository(dbRepository),
		Orders:   orders.NewRepository(dbRepository),
		Products: products.NewRepository(dbRepository),
	}
}
