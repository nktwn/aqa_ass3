package repository

import (
	"context"
	"diploma/modules/supplier/model"
	"diploma/modules/supplier/repo/converter"
	repoModel "diploma/modules/supplier/repo/model"
	"diploma/pkg/client/db"

	sq "github.com/Masterminds/squirrel"
)

func (r *supplierRepo) SupplierListByIDList(ctx context.Context, id []int64) ([]model.Supplier, error) {

	builder := sq.
		Select(
			"s."+sIDCol+" AS supplier_id",
			"s."+sNameCol+" AS supplier_name",
			"s."+sOrderAmountCol+" AS order_amount",

			"dc."+dcFreeDeliveryAmountCol+" AS minimum_free_delivery_amount",
			"dc."+dcDeliveryFeeCol+" AS delivery_fee",
		).
		From(supplierTbl + " AS s").
		LeftJoin(deliveryConditionTbl + " AS dc ON dc." + dcIDCol + " = s." + sDeliveryConditionIDCol).
		Where(sq.Eq{sIDCol: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "product_repository.GetSupplierProductListByProduct",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []repoModel.Supplier
	for rows.Next() {
		var s repoModel.Supplier
		if err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.OrderAmount,
			&s.FreeDeliveryAmount,
			&s.DeliveryFee,
		); err != nil {
			return nil, err
		}
		results = append(results, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return converter.ToServiceSupplierListFromRepo(results), nil

}
