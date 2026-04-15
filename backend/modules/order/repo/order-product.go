package repository

import (
	"context"
	"diploma/modules/order/model"
	"diploma/modules/order/repo/converter"
	modelRepo "diploma/modules/order/repo/model"
	"diploma/pkg/client/db"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

const (
	orderProductsTable = "order_products"
	opOrderIDColumn    = "order_id"
	opProductIDColumn  = "product_id"
	opQuantityColumn   = "quantity"
	opPriceColumn      = "price"
)

// CreateOrderProduct inserts a new order-product relationship.
func (r *OrderRepo) CreateOrderProduct(ctx context.Context, orderProduct *model.OrderProduct) error {
	builder := sq.
		Insert(orderProductsTable).
		PlaceholderFormat(sq.Dollar).
		Columns(opOrderIDColumn, opProductIDColumn, opQuantityColumn, opPriceColumn).
		Values(orderProduct.OrderID, orderProduct.ProductID, orderProduct.Quantity, orderProduct.Price)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "order_repository.CreateOrderProduct",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return nil
}

// GetOrderProducts retrieves all products for a specific order.
func (r *OrderRepo) OrderProducts(ctx context.Context, orderID int64) ([]*model.OrderProduct, error) {
	builder := sq.
		Select(opOrderIDColumn, opProductIDColumn, opQuantityColumn, opPriceColumn).
		PlaceholderFormat(sq.Dollar).
		From(orderProductsTable).
		Where(sq.Eq{opOrderIDColumn: orderID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "order_repository.GetOrderProducts",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderProducts []*model.OrderProduct
	for rows.Next() {
		var op modelRepo.OrderProduct
		err := rows.Scan(
			&op.OrderID,
			&op.ProductID,
			&op.Quantity,
			&op.Price,
		)
		if err != nil {
			return nil, err
		}
		orderProducts = append(orderProducts, converter.ToServiceOrderProductFromRepo(op))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orderProducts, nil
}

// UpdateOrderProduct modifies an existing order-product relationship.
func (r *OrderRepo) UpdateOrderProduct(ctx context.Context, orderProduct *model.OrderProduct) error {
	builder := sq.
		Update(orderProductsTable).
		PlaceholderFormat(sq.Dollar).
		Set(opQuantityColumn, orderProduct.Quantity).
		Set(opPriceColumn, orderProduct.Price).
		Where(sq.And{
			sq.Eq{opOrderIDColumn: orderProduct.OrderID},
			sq.Eq{opProductIDColumn: orderProduct.ProductID},
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "order_repository.UpdateOrderProduct",
		QueryRaw: query,
	}

	result, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no rows updated")
	}
	return nil
}

// DeleteOrderProduct removes an order-product relationship.
func (r *OrderRepo) DeleteOrderProduct(ctx context.Context, orderID, productID, supplierID int64) error {
	builder := sq.
		Delete(orderProductsTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.And{
			sq.Eq{opOrderIDColumn: orderID},
			sq.Eq{opProductIDColumn: productID},
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "order_repository.DeleteOrderProduct",
		QueryRaw: query,
	}

	result, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no rows deleted")
	}
	return nil
}
