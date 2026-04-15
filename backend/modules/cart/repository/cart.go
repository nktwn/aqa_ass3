package repository

import (
	"context"
	"database/sql"
	"diploma/modules/cart/model"
	"diploma/modules/cart/repository/converter"
	modelRepo "diploma/modules/cart/repository/model"
	"diploma/pkg/client/db"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

// CreateCart creates a new cart for the given user
func (r *cartRepo) CreateCart(ctx context.Context, userID int64) (int64, error) {
	builder := sq.Insert(cartsTable).
		PlaceholderFormat(sq.Dollar).
		Columns(cCustomerIDColumn).
		Values(userID).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build create cart query: %w", err)
	}

	q := db.Query{
		Name:     "cart_repository.Create",
		QueryRaw: query,
	}

	var cartID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&cartID)
	if err != nil {
		return 0, fmt.Errorf("failed to create cart: %w", err)
	}

	return cartID, nil
}

// Cart retrieves a cart for the given user
func (r *cartRepo) Cart(ctx context.Context, userID int64) (*model.Cart, error) {
	builder := sq.Select(cIDColumn, cTotalColumn, cCustomerIDColumn).
		PlaceholderFormat(sq.Dollar).
		From(cartsTable).
		Where(sq.Eq{cCustomerIDColumn: userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get cart query: %w", err)
	}

	q := db.Query{
		Name:     "cart_repository.Get",
		QueryRaw: query,
	}

	var cart modelRepo.Cart
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&cart.ID, &cart.Total, &cart.CustomerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return converter.ToServiceCartFromRepo(&cart), nil
}

// UpdateCartTotal updates the total amount for the given cart
func (r *cartRepo) UpdateCartTotal(ctx context.Context, cartID int64, total int) error {
	builder := sq.Update(cartsTable).
		PlaceholderFormat(sq.Dollar).
		Set(cTotalColumn, total).
		Where(sq.Eq{cIDColumn: cartID})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update cart total query: %w", err)
	}

	q := db.Query{
		Name:     "cart_repository.UpdateTotal",
		QueryRaw: query,
	}

	result, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to update cart total: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return model.ErrNoRows
	}

	return nil
}

// GetCartItems retrieves all items in the given cart with their product details
func (r *cartRepo) GetCartItems(ctx context.Context, cartID int64) ([]model.Supplier, error) {
	builder := sq.Select(
		cartItemTable+"."+ciProductIDColumn,
		cartItemTable+"."+ciSupplierIDColumn,
		cartItemTable+"."+ciQuantityColumn,
		cartItemTable+"."+ciPriceColumn,
		productsTable+"."+nameColumn+" AS product_name",
		productsTable+"."+imageURLColumn,
	).
		PlaceholderFormat(sq.Dollar).
		From(cartItemTable).
		Join(productsTable + " ON " + cartItemTable + "." + ciProductIDColumn + " = " + productsTable + ".id").
		Where(sq.Eq{ciCartIDColumn: cartID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get cart items query: %w", err)
	}

	q := db.Query{
		Name:     "cart_repository.GetCartItems",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}
	defer rows.Close()

	items := make([]modelRepo.CartItem, 0)

	for rows.Next() {
		var item modelRepo.CartItem
		if err := rows.Scan(
			&item.ProductID,
			&item.SupplierID,
			&item.Quantity,
			&item.Price,
			&item.ProductName,
			&item.ProductImageURL,
		); err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cart items: %w", err)
	}

	if len(items) == 0 {
		return []model.Supplier{}, nil
	}

	return converter.ToServiceSupplierFromRepo(items), nil
}

// DeleteCart removes the cart and all its items
func (r *cartRepo) DeleteCart(ctx context.Context, cartID int64) error {
	builder := sq.Delete(cartsTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{cIDColumn: cartID})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete cart query: %w", err)
	}

	q := db.Query{
		Name:     "cart_repository.Delete",
		QueryRaw: query,
	}

	result, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to delete cart: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return model.ErrNoRows
	}

	return nil
}
