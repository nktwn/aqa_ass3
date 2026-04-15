package redis

import (
	"context"
	"encoding/json"

	"diploma/modules/cart/model"
)

// SavePaymentOrder marshals the PaymentOrder into JSON and stores it in Redis.
func (cr *cartRedis) SavePaymentOrder(ctx context.Context, order model.PaymentOrder) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	key := "paymentorder:" + order.ID
	// Set expiration to 0 for no expiration. Adjust if needed.
	return cr.redisCl.Set(ctx, key, data, 0).Err()
}

// GetPaymentOrder retrieves a PaymentOrder from Redis by its ID and unmarshals it from JSON.
func (cr *cartRedis) PaymentOrder(ctx context.Context, id string) (model.PaymentOrder, error) {
	key := "paymentorder:" + id
	data, err := cr.redisCl.Get(ctx, key).Bytes()
	if err != nil {
		return model.PaymentOrder{}, err
	}

	var order model.PaymentOrder
	if err = json.Unmarshal(data, &order); err != nil {
		return model.PaymentOrder{}, err
	}
	return order, nil
}

// PaymentOrdersByPrefix retrieves all PaymentOrders from Redis whose IDs start with the given prefix.
// It builds a pattern matching keys formatted as "paymentorder:<prefix>*".
// Note: Using the KEYS command can be inefficient in production environments.
func (cr *cartRedis) PaymentOrdersByPrefix(ctx context.Context, prefix string) ([]*model.PaymentOrder, error) {
	pattern := "paymentorder:" + prefix + "*"
	keys, err := cr.redisCl.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var orders []*model.PaymentOrder
	for _, key := range keys {
		data, err := cr.redisCl.Get(ctx, key).Bytes()
		if err != nil {
			// continue if a key doesn't give valid data
			continue
		}
		var order model.PaymentOrder
		if err = json.Unmarshal(data, &order); err != nil {
			continue
		}
		orders = append(orders, &order)
	}
	return orders, nil
}
