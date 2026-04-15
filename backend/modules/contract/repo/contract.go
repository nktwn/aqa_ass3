package repo

import (
	"context"
	"diploma/modules/contract/model"
	"diploma/pkg/client/db"

	sq "github.com/Masterminds/squirrel"
)

type Repository struct {
	db db.Client
}

func NewRepository(db db.Client) *Repository {
	return &Repository{db: db}
}

const contractTable = "contracts"

func (r *Repository) Create(ctx context.Context, contract *model.Contract) (int64, error) {
	query, args, _ := sq.Insert(contractTable).
		Columns("order_id", "supplier_id", "customer_id", "content", "status", "created_at").
		Values(contract.OrderID, contract.SupplierID, contract.CustomerID, contract.Content, contract.Status, contract.CreatedAt).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	var id int64
	err := r.db.DB().QueryRowContext(ctx, db.Query{Name: "contract.create", QueryRaw: query}, args...).Scan(&id)
	return id, err
}

func (r *Repository) SignByParty(ctx context.Context, contractID int64, role int, signature string) error {
	update := sq.Update(contractTable).
		Where(sq.Eq{"id": contractID}).
		PlaceholderFormat(sq.Dollar)

	if role == 0 {
		update = update.Set("customer_sig", signature).Set("status", model.StatusSignedByCustomer)
	} else if role == 1 {
		update = update.Set("supplier_sig", signature).Set("status", model.StatusSignedBySupplier)
	}

	query, args, _ := update.ToSql()
	_, err := r.db.DB().ExecContext(ctx, db.Query{Name: "contract.sign", QueryRaw: query}, args...)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*model.Contract, error) {
	query, args, _ := sq.Select("*").From(contractTable).
		Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()

	row := r.db.DB().QueryRowContext(ctx, db.Query{Name: "contract.get", QueryRaw: query}, args...)

	var c model.Contract
	err := row.Scan(&c.ID, &c.OrderID, &c.SupplierID, &c.CustomerID, &c.Content, &c.SupplierSig, &c.CustomerSig, &c.Status, &c.CreatedAt, &c.SignedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) GetByUser(ctx context.Context, userID int64) ([]*model.Contract, error) {
	query, args, _ := sq.Select("*").
		From(contractTable).
		Where(sq.Or{
			sq.Eq{"supplier_id": userID},
			sq.Eq{"customer_id": userID},
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	rows, err := r.db.DB().QueryContext(ctx, db.Query{Name: "contract.get_by_user", QueryRaw: query}, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contracts []*model.Contract
	for rows.Next() {
		var c model.Contract
		err := rows.Scan(
			&c.ID, &c.OrderID, &c.SupplierID, &c.CustomerID,
			&c.Content, &c.SupplierSig, &c.CustomerSig,
			&c.Status, &c.CreatedAt, &c.SignedAt,
		)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contracts, nil
}

func (r *Repository) MarkAsSigned(ctx context.Context, contractID int64) error {
	query, args, _ := sq.Update(contractTable).
		Set("status", model.StatusCompleted).
		Set("signed_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": contractID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	_, err := r.db.DB().ExecContext(ctx, db.Query{
		Name:     "contract.mark_as_signed",
		QueryRaw: query,
	}, args...)
	return err
}
