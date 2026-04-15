package service

import (
	"context"
	"diploma/modules/order/model"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type IOrderRepo interface {
	OrdersByUserID(ctx context.Context, userID int64) ([]*model.Order, error)
	OrdersBySupplierID(ctx context.Context, supplierID int64) ([]*model.Order, error)
	OrderProducts(ctx context.Context, orderID int64) ([]*model.OrderProduct, error)
}

func (s *OrderService) Orders(ctx context.Context, userID int64, role int) ([]*model.Order, error) {
	s.LogInfo(ctx, "Fetching orders",
		zap.Int64("user_id", userID),
		zap.Int("role", role),
	)

	var orders []*model.Order
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		switch role {
		case model.CustomerRole:
			s.LogDebug(ctx, "Fetching customer orders")
			orders, errTx = s.orderRepo.OrdersByUserID(ctx, userID)
			if errTx != nil {
				s.LogError(ctx, "Failed to fetch customer orders", errTx)
				return errTx
			}
			for _, order := range orders {
				order.Supplier, errTx = s.supplierClient.Supplier(ctx, order.SupplierID)
				if errTx != nil {
					s.LogError(ctx, "Failed to fetch supplier details", errTx,
						zap.Int64("supplier_id", order.SupplierID),
					)
					return errTx
				}
				order.Supplier.ID = order.SupplierID
			}

		case model.SupplierRole:
			s.LogDebug(ctx, "Fetching supplier orders")
			orders, errTx = s.orderRepo.OrdersBySupplierID(ctx, userID)
			if errTx != nil {
				s.LogError(ctx, "Failed to fetch supplier orders", errTx)
				return errTx
			}
			for _, order := range orders {
				order.Supplier, errTx = s.supplierClient.Supplier(ctx, order.SupplierID)
				if errTx != nil {
					s.LogError(ctx, "Failed to fetch supplier details", errTx,
						zap.Int64("supplier_id", order.SupplierID),
					)
					return errTx
				}
			}
		}

		for _, order := range orders {
			orderProducts, errTx := s.orderRepo.OrderProducts(ctx, order.ID)
			if errTx != nil {
				s.LogError(ctx, "Failed to fetch order products", errTx,
					zap.Int64("order_id", order.ID),
				)
				return errTx
			}
			order.ProductList = orderProducts
			for _, orderProduct := range orderProducts {
				product, errTx := s.productClient.Product(ctx, orderProduct.ProductID)
				if errTx != nil {
					s.LogError(ctx, "Failed to fetch product details", errTx,
						zap.Int64("product_id", orderProduct.ProductID),
					)
					return errTx
				}
				orderProduct.Product = product
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	s.LogInfo(ctx, "Successfully fetched orders",
		zap.Int("order_count", len(orders)),
		zap.Int64("user_id", userID),
		zap.Int("role", role),
	)
	return orders, nil
}

type ICreateOrderRepo interface {
	CreateOrder(ctx context.Context, order *model.Order) (int64, error)
	CreateOrderProduct(ctx context.Context, orderProduct *model.OrderProduct) error
}

func (s *OrderService) CreateOrder(ctx context.Context, orders []*model.Order) error {
	s.LogInfo(ctx, "Creating orders", zap.Int("order_count", len(orders)))

	ordersID := make([]int64, len(orders))
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		for _, order := range orders {
			order.OrderDate = time.Now().Add(72 * time.Hour)
			order.StatusID = 1

			s.LogDebug(ctx, "Creating order",
				zap.Int64("customer_id", order.CustomerID),
				zap.Int64("supplier_id", order.SupplierID),
				zap.Time("order_date", order.OrderDate),
			)

			id, errTx := s.orderRepo.CreateOrder(ctx, order)
			if errTx != nil {
				s.LogError(ctx, "Failed to create order", errTx)
				return errTx
			}

			for _, op := range order.ProductList {
				op.OrderID = id
				s.LogDebug(ctx, "Adding product to order",
					zap.Int64("order_id", id),
					zap.Int64("product_id", op.ProductID),
					zap.Int("quantity", op.Quantity),
				)

				errTx = s.orderRepo.CreateOrderProduct(ctx, op)
				if errTx != nil {
					s.LogError(ctx, "Failed to create order product", errTx,
						zap.Int64("order_id", id),
						zap.Int64("product_id", op.ProductID),
					)
					return errTx
				}
			}
			ordersID = append(ordersID, id)

			if s.notificationService != nil {
				orderID := id
				if errTx := s.notificationService.Notify(ctx, order.CustomerID, "Заказ создан", fmt.Sprintf("Заказ #%d успешно оформлен.", id), "order_created", &orderID); errTx != nil {
					return errTx
				}
				if errTx := s.notificationService.Notify(ctx, order.SupplierID, "Новый заказ", fmt.Sprintf("Поступил новый заказ #%d.", id), "order_created", &orderID); errTx != nil {
					return errTx
				}
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	s.LogInfo(ctx, "Successfully created orders",
		zap.Int("order_count", len(orders)),
		zap.Any("order_ids", ordersID),
	)
	return nil
}

func (s *OrderService) UpdateOrderStatusBySupplier(ctx context.Context, supplierID int64, orderID int64, newStatusID int) error {
	s.LogInfo(ctx, "Updating order status",
		zap.Int64("supplier_id", supplierID),
		zap.Int64("order_id", orderID),
		zap.Int("new_status", newStatusID),
	)

	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		order, err := s.orderRepo.GetOrderByID(ctx, orderID)
		if err != nil {
			s.LogError(ctx, "Failed to get order", err,
				zap.Int64("order_id", orderID),
			)
			return err
		}

		if order.SupplierID != supplierID {
			s.LogWarn(ctx, "Unauthorized status update attempt",
				zap.Int64("supplier_id", supplierID),
				zap.Int64("order_id", orderID),
				zap.Int64("actual_supplier_id", order.SupplierID),
			)
			return fmt.Errorf("supplier %d does not own order %d", supplierID, orderID)
		}

		valid := false
		switch order.StatusID {
		case model.Pending:
			if newStatusID == model.InProgress {
				valid = true
			}
		case model.InProgress:
			if newStatusID == model.Completed || newStatusID == model.Cancelled {
				valid = true
			}
		}

		if !valid {
			s.LogWarn(ctx, "Invalid status transition",
				zap.Int("current_status", order.StatusID),
				zap.Int("new_status", newStatusID),
			)
			return fmt.Errorf("invalid status transition from %d to %d", order.StatusID, newStatusID)
		}

		err = s.orderRepo.UpdateOrderStatus(ctx, orderID, newStatusID)
		if err != nil {
			s.LogError(ctx, "Failed to update order status", err)
			return err
		}

		if order.StatusID == model.Pending && newStatusID == model.InProgress {
			s.LogInfo(ctx, "Creating contract for order",
				zap.Int64("order_id", order.ID),
			)
			_, err := s.contractClient.CreateContract(
				ctx,
				order.ID,
				order.SupplierID,
				order.CustomerID,
				fmt.Sprintf("Контракт для заказа #%d", order.ID),
			)
			if err != nil {
				s.LogError(ctx, "Failed to create contract", err)
				return err
			}
		}

		if s.notificationService != nil {
			orderIDPtr := order.ID
			err = s.notificationService.Notify(
				ctx,
				order.CustomerID,
				"Статус заказа обновлён",
				fmt.Sprintf("Заказ #%d теперь имеет статус: %s.", order.ID, statusName(newStatusID)),
				"order_status",
				&orderIDPtr,
			)
			if err != nil {
				s.LogError(ctx, "Failed to create notification", err)
				return err
			}
		}

		s.LogInfo(ctx, "Successfully updated order status",
			zap.Int64("order_id", orderID),
			zap.Int("old_status", order.StatusID),
			zap.Int("new_status", newStatusID),
		)
		return nil
	})
}

func (s *OrderService) GetOrderByID(ctx context.Context, userID int64, role int, orderID int64) (*model.Order, error) {
	s.LogInfo(ctx, "Fetching order details",
		zap.Int64("user_id", userID),
		zap.Int("role", role),
		zap.Int64("order_id", orderID),
	)

	var order *model.Order
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		order, errTx = s.orderRepo.GetOrderByID(ctx, orderID)
		if errTx != nil {
			s.LogError(ctx, "Failed to get order", errTx)
			return errTx
		}

		switch role {
		case model.CustomerRole:
			if order.CustomerID != userID {
				s.LogWarn(ctx, "Unauthorized order access attempt",
					zap.Int64("user_id", userID),
					zap.Int64("order_id", orderID),
					zap.Int64("order_customer_id", order.CustomerID),
				)
				return fmt.Errorf("unauthorized: customer can only view their own orders")
			}
		case model.SupplierRole:
			if order.SupplierID != userID {
				s.LogWarn(ctx, "Unauthorized order access attempt",
					zap.Int64("user_id", userID),
					zap.Int64("order_id", orderID),
					zap.Int64("order_supplier_id", order.SupplierID),
				)
				return fmt.Errorf("unauthorized: supplier can only view their own orders")
			}
		case model.AdminRole:
		default:
			s.LogWarn(ctx, "Invalid role access attempt",
				zap.Int("role", role),
			)
			return fmt.Errorf("unauthorized role")
		}

		order.Supplier, errTx = s.supplierClient.Supplier(ctx, order.SupplierID)
		if errTx != nil {
			s.LogError(ctx, "Failed to get supplier details", errTx)
			return errTx
		}

		orderProducts, errTx := s.orderRepo.OrderProducts(ctx, order.ID)
		if errTx != nil {
			s.LogError(ctx, "Failed to get order products", errTx)
			return errTx
		}
		order.ProductList = orderProducts

		for _, op := range order.ProductList {
			op.Product, errTx = s.productClient.Product(ctx, op.ProductID)
			if errTx != nil {
				s.LogError(ctx, "Failed to get product details", errTx,
					zap.Int64("product_id", op.ProductID),
				)
				return errTx
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.LogInfo(ctx, "Successfully fetched order details",
		zap.Int64("order_id", orderID),
		zap.Int("product_count", len(order.ProductList)),
	)
	return order, nil
}

func (s *OrderService) CancelOrderByCustomer(ctx context.Context, customerID int64, orderID int64) error {
	s.LogInfo(ctx, "Cancelling order",
		zap.Int64("customer_id", customerID),
		zap.Int64("order_id", orderID),
	)

	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		order, err := s.orderRepo.GetOrderByID(ctx, orderID)
		if err != nil {
			s.LogError(ctx, "Failed to get order", err)
			return err
		}

		if order.CustomerID != customerID {
			s.LogWarn(ctx, "Unauthorized cancellation attempt",
				zap.Int64("customer_id", customerID),
				zap.Int64("order_id", orderID),
				zap.Int64("order_customer_id", order.CustomerID),
			)
			return fmt.Errorf("customer %d does not own order %d", customerID, orderID)
		}

		if order.StatusID != model.Pending {
			s.LogWarn(ctx, "Invalid order status for cancellation",
				zap.Int("current_status", order.StatusID),
			)
			return fmt.Errorf("only orders in Pending status can be cancelled")
		}

		err = s.orderRepo.UpdateOrderStatus(ctx, orderID, model.Cancelled)
		if err != nil {
			s.LogError(ctx, "Failed to update order status", err)
			return err
		}

		if s.notificationService != nil {
			orderIDPtr := order.ID
			if err := s.notificationService.Notify(
				ctx,
				order.SupplierID,
				"Заказ отменён",
				fmt.Sprintf("Покупатель отменил заказ #%d.", order.ID),
				"order_cancelled",
				&orderIDPtr,
			); err != nil {
				return err
			}
		}

		s.LogInfo(ctx, "Successfully cancelled order",
			zap.Int64("order_id", orderID),
		)
		return nil
	})
}

func statusName(statusID int) string {
	switch statusID {
	case model.Pending:
		return "Pending"
	case model.InProgress:
		return "In Progress"
	case model.Completed:
		return "Completed"
	case model.Cancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}
