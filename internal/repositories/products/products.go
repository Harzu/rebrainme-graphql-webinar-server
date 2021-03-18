package products

import (
	"context"
	"fmt"

	"webinar/graphql/server/internal/entities"
)

func (r *repositoryDB) InsertOrUpdateProduct(_ context.Context, name string, price int64) (int64, error) {
	query, args, err := prepareInsertOrUpdateProduct(name, price)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare InsertOrUpdateProduct query: %s", err)
	}

	var userId int64
	returningRow := r.client.QueryRowx(query, args...)
	if err := returningRow.Scan(&userId); err != nil {
		return 0, fmt.Errorf("failed to scan returning row to productId: %w", err)
	}

	return userId, nil
}

func (r *repositoryDB) FindProductsByIds(_ context.Context, productIds []int64) ([]entities.Product, error) {
	query, args, err := prepareFindProductsByIdsQuery(productIds)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare FindProductsByIds query: %w", err)
	}

	rows, err := r.client.Queryx(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute FindProductsByIds: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var result []entities.Product
	for rows.Next() {
		var product entities.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			return nil, fmt.Errorf("failed to scan row to Product model: %w", err)
		}
		result = append(result, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("unable to scan all out of FindProductsByIds: %w", err)
	}

	return result, err
}

func (r *repositoryDB) FindProductsByOrdersMap(_ context.Context, orderIds []int64) (_ map[int64][]entities.Product, err error) {
	result := map[int64][]entities.Product{}

	query, args, err := prepareFindProductsByOrdersMapQuery(orderIds)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare FindProductsByOrdersMap query: %w", err)
	}

	rows, err := r.client.Queryx(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute FindProductsByOrdersMap: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	for rows.Next() {
		var model orderProductModel
		if err := rows.StructScan(&model); err != nil {
			return nil, fmt.Errorf("failed to scan row to OrderProduct: %w", err)
		}

		if _, exists := result[model.OrderID]; !exists {
			result[model.OrderID] = []entities.Product{}
		}

		result[model.OrderID] = append(result[model.OrderID], buildProductEntity(model))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("unable to scan all out of FindProductsByOrdersMap: %w", err)
	}

	return result, err
}
