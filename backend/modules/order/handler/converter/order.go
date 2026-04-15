package converter

import (
	apiModel "diploma/modules/order/handler/model"
	serviceModel "diploma/modules/order/model"
)

// ConvertOrderToAPI converts service Order model to API Order response
func ConvertOrderToAPI(order *serviceModel.Order) *apiModel.Order {
	return &apiModel.Order{
		ID:              order.ID,
		Status:          convertStatusIDToString(order.StatusID),
		OrderDate:       order.OrderDate.Format("2006-01-02T15:04:05Z07:00"),
		DeliveryAddress: order.DeliveryAddress,
		DeliveryComment: order.DeliveryComment,
		PaymentStatus:   order.PaymentStatus,
		Supplier: &apiModel.Supplier{
			ID:   order.Supplier.ID,
			Name: order.Supplier.Name,
		},
		ProductList: convertOrderProductsToAPI(order.ProductList),
	}
}

// convertOrderProductsToAPI converts a slice of OrderProduct to a slice of Product for API response
func convertOrderProductsToAPI(products []*serviceModel.OrderProduct) []*apiModel.Product {
	apiProducts := make([]*apiModel.Product, len(products))
	for i, p := range products {
		apiProducts[i] = &apiModel.Product{
			ID:       p.Product.ID,
			Name:     p.Product.Name,
			ImageUrl: p.Product.ImageUrl,
			Quantity: p.Quantity,
			Price:    p.Price,
		}
	}
	return apiProducts
}

// convertStatusIDToString converts a status ID to a string representation
func convertStatusIDToString(statusID int) string {
	switch statusID {
	case serviceModel.Pending:
		return "Pending"
	case serviceModel.InProgress:
		return "In Progress"
	case serviceModel.Completed:
		return "Completed"
	case serviceModel.Cancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}

// ConvertOrdersToAPI converts multiple service Order models to an API response
func ConvertOrdersToAPI(orders []*serviceModel.Order) *apiModel.GetOrdersResponse {
	apiOrders := make([]apiModel.Order, len(orders))
	for i, order := range orders {
		apiOrders[i] = *ConvertOrderToAPI(order)
	}
	return &apiModel.GetOrdersResponse{
		Orders: apiOrders,
	}
}
