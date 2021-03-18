package customers

import (
	"context"
	"errors"

	"webinar/graphql/server/internal/constants"
	"webinar/graphql/server/internal/entities"
	"webinar/graphql/server/internal/graph/model"
	ordersRepo "webinar/graphql/server/internal/repositories/orders"
	"webinar/graphql/server/internal/repositories/products"
	"webinar/graphql/server/internal/repositories/users"
)

type CustomerResolvers interface {
	Me(ctx context.Context) (*model.Customer, error)
}

type resolvers struct {
	usersRepo    users.Repository
	ordersRepo   ordersRepo.Repository
	productsRepo products.Repository
}

func New(
	usersRepo users.Repository,
	ordersRepo ordersRepo.Repository,
	productsRepo products.Repository,
) CustomerResolvers {
	return &resolvers{
		usersRepo:    usersRepo,
		ordersRepo:   ordersRepo,
		productsRepo: productsRepo,
	}
}

func (r *resolvers) Me(ctx context.Context) (*model.Customer, error) {
	session, ok := ctx.Value(constants.SessionContextKey).(entities.Session)
	if !ok {
		return nil, errors.New("failed to get user session from context")
	}

	customer, err := r.usersRepo.FindCustomerByUserId(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	orders, err := r.ordersRepo.FindOrdersByCustomerId(ctx, customer.ID)
	if err != nil {
		return nil, err
	}

	productsByOrderMap, err := r.productsRepo.FindProductsByOrdersMap(ctx, orders.ExtractIds())
	if err != nil {
		return nil, err
	}

	ordersModel := make([]*model.Order, 0, len(orders))
	for index := range orders {
		orderModel := &model.Order{
			ID:         orders[index].ID,
			CustomerID: orders[index].CustomerID,
			TotalPrice: orders[index].TotalPrice,
			Status:     orders[index].Status,
		}

		orderProducts := productsByOrderMap[orders[index].ID]
		orderModel.Products = make([]*model.Product, 0, len(orderProducts))
		for orderProductIndex := range orderProducts {
			orderProductModel := &model.Product{
				ID:    orderProducts[orderProductIndex].ID,
				Name:  orderProducts[orderProductIndex].Name,
				Price: orderProducts[orderProductIndex].Price,
			}
			orderModel.Products = append(orderModel.Products, orderProductModel)
		}
		ordersModel = append(ordersModel, orderModel)
	}

	return &model.Customer{
		ID:      customer.ID,
		UserID:  customer.UserID,
		Name:    customer.Name,
		Address: customer.Address,
		Orders:  ordersModel,
	}, nil
}
