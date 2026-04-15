package repository

import "diploma/pkg/client/db"

type supplierRepo struct {
	db db.Client
}

const (
	// ======== supplier table ========
	supplierTbl             = "suppliers"
	sIDCol                  = "user_id"
	sNameCol                = "name"
	sOrderAmountCol         = "order_amount"
	sDeliveryConditionIDCol = "condition_id"

	// ======== delivery conditions ========
	deliveryConditionTbl    = "delivery_conditions"
	dcIDCol                 = "condition_id"
	dcFreeDeliveryAmountCol = "minimum_free_delivery_amount"
	dcDeliveryFeeCol        = "delivery_fee"
)

func NewRepository(db db.Client) *supplierRepo {
	return &supplierRepo{db: db}
}
