package repository

import (
	"context"
	"database/sql"
	"diploma/modules/cart/model"
	"diploma/pkg/client/db"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

// ItemQuantity returns the quantity of a specific item in the cart
func (r *cartRepo) ItemQuantity(ctx context.Context, cartID, productId, supplierId int64) (int, error) {
	builder := sq.Select(ciQuantityColumn).
		PlaceholderFormat(sq.Dollar).
		From(cartItemTable).
		Where(sq.And{
			sq.Eq{ciCartIDColumn: cartID},
			sq.Eq{ciSupplierIDColumn: supplierId},
			sq.Eq{ciProductIDColumn: productId},
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build item quantity query: %w", err)
	}

	q := db.Query{
		Name:     "cart_repository.ItemQuantity",
		QueryRaw: query,
	}

	var quantity int
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return 0, model.ErrNoRows
		}
		return 0, fmt.Errorf("failed to get item quantity: %w", err)
	}

	return quantity, nil
}

// UpdateItemQuantity updates the quantity of a specific item in the cart
func (r *cartRepo) UpdateItemQuantity(ctx context.Context, cartID, productId, supplierId int64, quantity int) error {
	builder := sq.Update(cartItemTable).
		PlaceholderFormat(sq.Dollar).
		Set(ciQuantityColumn, quantity).
		Where(sq.And{
			sq.Eq{ciCartIDColumn: cartID},
			sq.Eq{ciSupplierIDColumn: supplierId},
			sq.Eq{ciProductIDColumn: productId},
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update item quantity query: %w", err)
	}

	q := db.Query{
		Name:     "cart_repository.UpdateItemQuantity",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to update item quantity: %w", err)
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return model.ErrNoRows
	}

	return nil
}

// AddItem adds a new item to the cart
func (r *cartRepo) AddItem(ctx context.Context, input *model.PutCartQuery) error {
	builder := sq.Insert(cartItemTable).
		PlaceholderFormat(sq.Dollar).
		Columns(ciCartIDColumn, ciSupplierIDColumn, ciProductIDColumn, ciPriceColumn, ciQuantityColumn).
		Values(input.CartID, input.SupplierID, input.ProductID, input.Price, input.Quantity)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build add item query: %w", err)
	}

	q := db.Query{
		Name:     "cart_repository.AddItem",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to add item: %w", err)
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return model.ErrNoRows
	}

	return nil
}

// DeleteCartItems removes all items from the cart
func (r *cartRepo) DeleteCartItems(ctx context.Context, cartID int64) error {
	builder := sq.Delete(cartItemTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{ciCartIDColumn: cartID})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete cart items query: %w", err)
	}

	q := db.Query{
		Name:     "cart_repository.DeleteCartItems",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to delete cart items: %w", err)
	}

	return nil
}

// DeleteItem removes a specific item from the cart
func (r *cartRepo) DeleteItem(ctx context.Context, cartID, productId, supplierId int64) error {
	builder := sq.Delete(cartItemTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.And{
			sq.Eq{ciCartIDColumn: cartID},
			sq.Eq{ciSupplierIDColumn: supplierId},
			sq.Eq{ciProductIDColumn: productId},
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete item query: %w", err)
	}

	q := db.Query{
		Name:     "cart_repository.DeleteItem",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return model.ErrNoRows
	}

	return nil
}
