package repository

import (
	"context"
	"diploma/modules/order/model"
	"diploma/modules/order/repo/converter"
	modelRepo "diploma/modules/order/repo/model"
	"diploma/pkg/client/db"

	sq "github.com/Masterminds/squirrel"
)

const (
	ordersTable       = "orders"
	oIDColumn         = "id"
	oCustomerIDColumn = "customer_id"
	oOrderDateColumn  = "order_date"
	oSupplierIDColumn = "supplier_id"
	oStatusIDColumn   = "status_id"
	oDeliveryAddress  = "delivery_address"
	oDeliveryComment  = "delivery_comment"
	oPaymentID        = "payment_id"
	oPaymentStatus    = "payment_status"
)

// CreateOrder inserts a new order record and returns its id.
func (r *OrderRepo) CreateOrder(ctx context.Context, order *model.Order) (int64, error) {
	builder := sq.
		Insert(ordersTable).
		PlaceholderFormat(sq.Dollar).
		Columns(oCustomerIDColumn, oSupplierIDColumn, oOrderDateColumn, oStatusIDColumn, oDeliveryAddress, oDeliveryComment, oPaymentID, oPaymentStatus).
		Values(order.CustomerID, order.SupplierID, order.OrderDate, order.StatusID, order.DeliveryAddress, order.DeliveryComment, order.PaymentID, order.PaymentStatus).
		Suffix("RETURNING " + oIDColumn)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "order_repository.CreateOrder",
		QueryRaw: query,
	}

	var orderID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&orderID)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

// GetOrder retrieves an order by its id.
func (r *OrderRepo) OrdersByUserID(ctx context.Context, userID int64) ([]*model.Order, error) {
	builder := sq.
		Select(oIDColumn, oCustomerIDColumn, oOrderDateColumn, oSupplierIDColumn, oStatusIDColumn, "COALESCE("+oDeliveryAddress+", '')", "COALESCE("+oDeliveryComment+", '')", "COALESCE("+oPaymentID+", '')", "COALESCE("+oPaymentStatus+", 'pending')").
		PlaceholderFormat(sq.Dollar).
		From(ordersTable).
		Where(sq.Eq{oCustomerIDColumn: userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "order_repository.OrdersByUserID",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		var orderRepo modelRepo.Order
		if err := rows.Scan(
			&orderRepo.ID,
			&orderRepo.CustomerID,
			&orderRepo.OrderDate,
			&orderRepo.SupplierID,
			&orderRepo.StatusID,
			&orderRepo.DeliveryAddress,
			&orderRepo.DeliveryComment,
			&orderRepo.PaymentID,
			&orderRepo.PaymentStatus,
		); err != nil {
			return nil, err
		}
		orders = append(orders, converter.ToServiceOrderFromRepo(&orderRepo))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return []*model.Order{}, nil
	}
	return orders, nil
}

// GetOrdersBySupplierID retrieves a list of orders by supplier id.
func (r *OrderRepo) OrdersBySupplierID(ctx context.Context, supplierID int64) ([]*model.Order, error) {
	builder := sq.
		Select(oIDColumn, oCustomerIDColumn, oOrderDateColumn, oSupplierIDColumn, oStatusIDColumn, "COALESCE("+oDeliveryAddress+", '')", "COALESCE("+oDeliveryComment+", '')", "COALESCE("+oPaymentID+", '')", "COALESCE("+oPaymentStatus+", 'pending')").
		PlaceholderFormat(sq.Dollar).
		From(ordersTable).
		Where(sq.Eq{oSupplierIDColumn: supplierID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "order_repository.GetOrdersBySupplierID",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		var orderRepo modelRepo.Order
		if err := rows.Scan(
			&orderRepo.ID,
			&orderRepo.CustomerID,
			&orderRepo.OrderDate,
			&orderRepo.SupplierID,
			&orderRepo.StatusID,
			&orderRepo.DeliveryAddress,
			&orderRepo.DeliveryComment,
			&orderRepo.PaymentID,
			&orderRepo.PaymentStatus,
		); err != nil {
			return nil, err
		}
		orders = append(orders, converter.ToServiceOrderFromRepo(&orderRepo))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return []*model.Order{}, nil
	}
	return orders, nil
}

func (r *OrderRepo) GetOrderByID(ctx context.Context, orderID int64) (*model.Order, error) {
	builder := sq.
		Select(oIDColumn, oCustomerIDColumn, oOrderDateColumn, oSupplierIDColumn, oStatusIDColumn, "COALESCE("+oDeliveryAddress+", '')", "COALESCE("+oDeliveryComment+", '')", "COALESCE("+oPaymentID+", '')", "COALESCE("+oPaymentStatus+", 'pending')").
		From(ordersTable).
		Where(sq.Eq{oIDColumn: orderID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "order_repository.GetOrderByID",
		QueryRaw: query,
	}

	var order modelRepo.Order
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&order.ID,
		&order.CustomerID,
		&order.OrderDate,
		&order.SupplierID,
		&order.StatusID,
		&order.DeliveryAddress,
		&order.DeliveryComment,
		&order.PaymentID,
		&order.PaymentStatus,
	)
	if err != nil {
		return nil, err
	}

	return converter.ToServiceOrderFromRepo(&order), nil
}

func (r *OrderRepo) UpdateOrderStatus(ctx context.Context, orderID int64, newStatus int) error {
	builder := sq.
		Update(ordersTable).
		Set(oStatusIDColumn, newStatus).
		Where(sq.Eq{oIDColumn: orderID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "order_repository.UpdateOrderStatus",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}

// UpdateOrder modifies an existing order record.
// func (r *OrderRepo) UpdateOrder(ctx context.Context, order *model.Order) error {
// 	builder := sq.
// 		Update(ordersTable).
// 		PlaceholderFormat(sq.Dollar).
// 		Set(oCustomerIDColumn, order.CustomerID).
// 		Set(oOrderDateColumn, order.OrderDate).

// 	query, args, err := builder.ToSql()
// 	if err != nil {
// 		return err
// 	}

// 	q := db.Query{
// 		Name:     "order_repository.UpdateOrder",
// 		QueryRaw: query,
// 	}

// 	result, err := r.db.DB().ExecContext(ctx, q, args...)
// 	if err != nil {
// 		return err
// 	}
// 	if affected, _ := result.RowsAffected(); affected == 0 {
// 		return fmt.Errorf("no rows updated")
// 	}
// 	return nil
// }

// DeleteOrder removes an order record by its id.
// func (r *OrderRepo) DeleteOrder(ctx context.Context, id int64) error {
// 	builder := sq.
// 		Delete(ordersTable).
// 		PlaceholderFormat(sq.Dollar).
// 		Where(sq.Eq{oIDColumn: id})

// 	query, args, err := builder.ToSql()
// 	if err != nil {
// 		return err
// 	}

// 	q := db.Query{
// 		Name:     "order_repository.DeleteOrder",
// 		QueryRaw: query,
// 	}

// 	result, err := r.db.DB().ExecContext(ctx, q, args...)
// 	if err != nil {
// 		return err
// 	}
// 	if affected, _ := result.RowsAffected(); affected == 0 {
// 		return fmt.Errorf("no rows deleted")
// 	}
// 	return nil
// }
