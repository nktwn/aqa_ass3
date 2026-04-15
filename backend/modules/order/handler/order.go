package handler

import (
	"context"
	"diploma/modules/auth/jwt"
	"diploma/modules/order/handler/converter"
	modelApi "diploma/modules/order/handler/model"
	"diploma/modules/order/model"
	contextkeys "diploma/pkg/context-keys"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ICreateOrderService interface {
	Orders(ctx context.Context, userID int64, role int) ([]*model.Order, error)
	// CreateOrder(userID int64) error
}

// GetOrders godoc
// @Summary Retrieve orders for a user
// @Description Retrieves orders for the authenticated user using the provided JWT claims.
// @Tags orders
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {array} modelApi.GetOrdersResponse "List of orders"
// @Failure 401 {object} modelApi.ErrorResponse "Unauthorized: invalid or missing JWT token"
// @Failure 500 {object} modelApi.ErrorResponse "Internal server error while retrieving orders"
// @Router /api/order [get]
func (h *OrderHandler) GetOrders(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
		return
	}

	// Retrieve the orders for the authenticated user
	orders, err := h.service.Orders(c.Request.Context(), claims.UserID, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, converter.ConvertOrdersToAPI(orders))
}

// GetOrderByID godoc
// @Summary Get order by ID
// @Description Retrieves an order by its ID.
// @Tags orders
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} modelApi.Order "Order details"
// @Failure 400 {object} modelApi.ErrorResponse "Invalid order ID"
// @Failure 401 {object} modelApi.ErrorResponse "Unauthorized: invalid or missing JWT token"
// @Failure 500 {object} modelApi.ErrorResponse "Internal server error while retrieving order"
// @Router /api/order/{id} [get]
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
		return
	}

	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "invalid order ID"})
		return
	}

	order, err := h.service.GetOrderByID(c.Request.Context(), claims.UserID, claims.Role, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, converter.ConvertOrderToAPI(order))
}

// UpdateOrderStatus godoc
// @Summary Update order status by supplier
// @Description Supplier updates the status of their order
// @Tags orders
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body modelApi.UpdateOrderStatusRequest true "Order ID and New Status"
// @Success 200 {object} map[string]string "status updated"
// @Failure 400 {object} modelApi.ErrorResponse "Invalid input"
// @Failure 401 {object} modelApi.ErrorResponse "Unauthorized"
// @Failure 500 {object} modelApi.ErrorResponse "Internal error"
// @Router /api/order/status [post]
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
		return
	}

	if claims.Role != model.SupplierRole {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: "only suppliers can update orders"})
		return
	}

	var req modelApi.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	err := h.service.UpdateOrderStatusBySupplier(c.Request.Context(), claims.UserID, req.OrderID, req.NewStatusID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order status updated"})
}

// CancelOrder godoc
// @Summary Cancel order by customer
// @Description Allows a customer to cancel their own order only if it's in Pending status.
// @Tags orders
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body modelApi.CancelOrderRequest true "Order ID to cancel"
// @Success 200 {object} map[string]string "order cancelled"
// @Failure 400 {object} modelApi.ErrorResponse "Invalid input or bad request"
// @Failure 401 {object} modelApi.ErrorResponse "Unauthorized"
// @Failure 403 {object} modelApi.ErrorResponse "Forbidden"
// @Failure 500 {object} modelApi.ErrorResponse "Internal server error"
// @Router /api/order/cancel [post]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
		return
	}

	if claims.Role != model.CustomerRole {
		c.JSON(http.StatusForbidden, modelApi.ErrorResponse{Err: "only customers can cancel orders"})
		return
	}

	var req modelApi.CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	err := h.service.CancelOrderByCustomer(c.Request.Context(), claims.UserID, req.OrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order cancelled"})
}

// CreateOrder is currently disabled and not registered in routes.
// func (h *OrderHandler) CreateOrder(c *gin.Context) {
// 	// Extract the order details from the request
// 	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)

// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
// 		return
// 	}
// 	err := h.service.CreateOrder(claims.UserID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"message": "Order created successfully",
// 	})
// }
