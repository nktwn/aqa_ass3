package repository

import (
	"context"
	"diploma/pkg/client/db"
	"time"

	model "diploma/modules/user/model"

	sq "github.com/Masterminds/squirrel"
)

const (
	aTableName         = "address"
	aIdColumn          = "id"
	aStreetColumn      = "street"
	aDescriptionColumn = "description"
	aUserIDColumn      = "user_id"
)

func (r *repo) CreateAddress(ctx context.Context, addr model.Address) (int64, error) {
	builder := sq.Insert(aTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(aStreetColumn, aDescriptionColumn, aUserIDColumn).
		Values(addr.Street, addr.Description, addr.UserID).
		Suffix("RETURNING " + aIdColumn)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "address_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) AddressByUserId(ctx context.Context, userID int64) ([]model.Address, error) {
	builder := sq.Select(aIdColumn, aStreetColumn, aDescriptionColumn, aUserIDColumn).
		PlaceholderFormat(sq.Dollar).
		From(aTableName).
		Where(sq.Eq{aUserIDColumn: userID})
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "address_repository.GetByUserId",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []model.Address
	for rows.Next() {
		var addr model.Address
		if err = rows.Scan(&addr.ID, &addr.Street, &addr.Description, &addr.UserID); err != nil {
			return nil, err
		}
		addresses = append(addresses, addr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}

func (r *repo) GetById(ctx context.Context, id int64) (*model.Address, error) {
	builder := sq.Select(aIdColumn, aStreetColumn, aDescriptionColumn, aUserIDColumn).
		PlaceholderFormat(sq.Dollar).
		From(aTableName).
		Where(sq.Eq{aIdColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "address_repository.GetById",
		QueryRaw: query,
	}

	var addr model.Address
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&addr.ID, &addr.Street, &addr.Description, &addr.UserID)
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

func (r *repo) Update(ctx context.Context, id int64, street, description string) error {
	builder := sq.Update(aTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{aIdColumn: id}).
		Set(aStreetColumn, street).
		Set(aDescriptionColumn, description).
		Set("updated_at", time.Now()) // Optional: remove if table doesn't have updated_at column

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "address_repository.Update",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(aTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{aIdColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "address_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}
