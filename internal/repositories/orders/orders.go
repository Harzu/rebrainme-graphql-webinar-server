package orders

import (
	"context"
	"fmt"

	"webinar/graphql/server/internal/entities"
	"webinar/graphql/server/internal/repositories/common"
)

func (r *repositoryDB) InsertOrder(ctx context.Context, customerID int64, productIds []int64, totalPrice int64) (_ int64, err error) {
	tx, err := r.client.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin InsertOrder tx: %w", err)
	}
	defer func() {
		if commitErr := common.CommitTransaction(ctx, tx, err, "InsertOrder"); commitErr != nil && err == nil {
			err = commitErr
		}
	}()

	insertOrderQuery, args, err := prepareInsertOrderQuery(customerID, totalPrice)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare InsertOrder query: %w", err)
	}

	var orderId int64
	returningRow := tx.QueryRowx(insertOrderQuery, args...)
	if err := returningRow.Scan(&orderId); err != nil {
		return 0, fmt.Errorf("failed to scan returning row to orderId: %w", err)
	}

	insertOrderProductsQuery, args, err := prepareInsertOrderProductsQuery(orderId, productIds)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare InsertOrderProducts query: %w", err)
	}

	if _, err := tx.Exec(insertOrderProductsQuery, args...); err != nil {
		return 0, fmt.Errorf("failed to execute InsertOrderProducts query: %w", err)
	}

	return orderId, err
}

func (r *repositoryDB) FindOrdersByCustomerId(_ context.Context, customerId int64) (_ entities.Orders, err error) {
	query, args, err := prepareFindOrdersByCustomerIdQuery(customerId)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare FindOrdersByCustomerId query: %w", err)
	}

	rows, err := r.client.Queryx(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute FindOrdersByCustomerId query: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var orders entities.Orders
	for rows.Next() {
		var model ordersModel
		if err := rows.StructScan(&model); err != nil {
			return nil, fmt.Errorf("failed to scan row to order: %w", err)
		}
		orders = append(orders, buildOrderEntity(model))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("unable to scan all out of FindOrdersByCustomerId: %w", err)
	}

	return orders, err
}
