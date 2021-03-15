package products

import (
	"context"
	"fmt"

	"webinar/graphql/server/internal/graph/model"
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

func (r *repositoryDB) FindProductsByIds(_ context.Context, productIds []int64) ([]model.Product, error) {
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

	var result []model.Product
	for rows.Next() {
		var product model.Product
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
