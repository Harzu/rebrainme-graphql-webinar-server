package products

import (
	"context"
	"fmt"

	"webinar/graphql/server/internal/graph/model"
	"webinar/graphql/server/internal/repositories/products"
)

type ProductResolvers interface {
	CreateOneProduct(ctx context.Context, input model.CreateProductInput) (*model.Product, error)
}

type resolvers struct {
	productsRepo products.Repository
}

func New(productsRepo products.Repository) ProductResolvers {
	return &resolvers{
		productsRepo: productsRepo,
	}
}

func (r *resolvers) CreateOneProduct(ctx context.Context, input model.CreateProductInput) (*model.Product, error) {
	productId, err := r.productsRepo.InsertOrUpdateProduct(ctx, input.Name, input.Price)
	if err != nil {
		return nil, fmt.Errorf("failed to insert or update product: %w", err)
	}

	return &model.Product{
		ID:    productId,
		Name:  input.Name,
		Price: input.Price,
	}, nil
}
