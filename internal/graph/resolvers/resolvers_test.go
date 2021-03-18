package resolvers

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"webinar/graphql/server/internal/constants"
	"webinar/graphql/server/internal/entities"
	"webinar/graphql/server/internal/graph/generated"
	"webinar/graphql/server/internal/graph/model"
	"webinar/graphql/server/internal/repositories"
	"webinar/graphql/server/internal/services"
	"webinar/graphql/server/test/mocks/ordersRepoMocks"
	"webinar/graphql/server/test/mocks/productsRepoMocks"
	"webinar/graphql/server/test/mocks/usersRepoMocks"
)

func TestResolvers_Me(t *testing.T) {
	req := require.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	any := gomock.Any()

	usersRepo := usersRepoMocks.NewMockRepository(mockCtrl)
	ordersRepo := ordersRepoMocks.NewMockRepository(mockCtrl)
	productsRepo := productsRepoMocks.NewMockRepository(mockCtrl)

	serviceContainer := services.Container{}
	repoContainer := repositories.Container{
		Users:    usersRepo,
		Orders:   ordersRepo,
		Products: productsRepo,
	}

	gqlConfig := generated.Config{
		Resolvers: New(&serviceContainer, &repoContainer),
		Directives: generated.DirectiveRoot{
			Auth: func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []model.Role) (res interface{}, err error) {
				return next(
					context.WithValue(
						ctx,
						constants.SessionContextKey,
						entities.Session{UserID: 1, Role: model.RoleCustomer},
					),
				)
			},
		},
	}
	testGqlServer := handler.NewDefaultServer(generated.NewExecutableSchema(gqlConfig))
	c := client.New(testGqlServer)

	usersRepo.EXPECT().FindCustomerByUserId(any, int64(1)).Return(entities.Customer{
		ID:      1,
		UserID:  1,
		Name:    "Ivan",
		Address: "test_address",
	}, nil).Times(1)

	ordersRepo.EXPECT().FindOrdersByCustomerId(any, int64(1)).Return(entities.Orders{
		{
			ID:         1,
			CustomerID: 1,
			Status:     model.OrderStatusCreated,
			TotalPrice: 10,
		},
	}, nil).Times(1)

	productsRepo.EXPECT().FindProductsByOrdersMap(any, []int64{1}).Return(map[int64][]entities.Product{
		1: {
			{ID: 1, Name: "product_1", Price: 5},
			{ID: 2, Name: "product_2", Price: 5},
		},
	}, nil).Times(1)

	var resp struct{ Me model.Customer }
	err := c.Post(`{
		me {
			id
			name
			address
			orders {
				id
				products {
					id
				}
			}
		}
	}`, &resp)

	req.NoError(err)
	req.EqualValues(model.Customer{
		ID:      1,
		Name:    "Ivan",
		Address: "test_address",
		Orders: []*model.Order{
			{
				ID: 1,
				Products: []*model.Product{
					{ID: 1},
					{ID: 2},
				},
			},
		},
	}, resp.Me)
}
