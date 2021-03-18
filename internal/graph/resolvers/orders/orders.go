package orders

import (
	"context"

	"webinar/graphql/server/internal/graph/model"
	"webinar/graphql/server/internal/repositories/orders"
	"webinar/graphql/server/internal/repositories/products"
)

type OrderResolvers interface {
	CreateOneOrder(ctx context.Context, customerID int64, productIds []int64) (*model.Order, error)
}

type resolvers struct {
	ordersRepo   orders.Repository
	productsRepo products.Repository
}

func New(ordersRepo orders.Repository, productsRepo products.Repository) OrderResolvers {
	return &resolvers{
		ordersRepo:   ordersRepo,
		productsRepo: productsRepo,
	}
}

func (r *resolvers) CreateOneOrder(ctx context.Context, customerID int64, productIds []int64) (*model.Order, error) {
	productsSlice, err := r.productsRepo.FindProductsByIds(ctx, productIds)
	if err != nil {
		return nil, err
	}

	var totalPrice int64
	for _, product := range productsSlice {
		totalPrice += product.Price
	}

	orderId, err := r.ordersRepo.InsertOrder(ctx, customerID, productIds, totalPrice)
	if err != nil {
		return nil, err
	}

	result := &model.Order{
		ID:         orderId,
		CustomerID: customerID,
		Status:     model.OrderStatusCreated,
		TotalPrice: totalPrice,
		Products:   make([]*model.Product, 0, len(productsSlice)),
	}

	for index := range productsSlice {
		productModel := &model.Product{
			ID:    productsSlice[index].ID,
			Name:  productsSlice[index].Name,
			Price: productsSlice[index].Price,
		}
		result.Products = append(result.Products, productModel)
	}

	return result, nil
}
