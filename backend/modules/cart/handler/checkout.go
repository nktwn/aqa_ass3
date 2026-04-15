package handler

import (
	"diploma/modules/auth/jwt"
	"diploma/modules/cart/handler/converter"
	modelApi "diploma/modules/cart/handler/model"
	"diploma/modules/cart/model"
	contextkeys "diploma/pkg/context-keys"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Checkout godoc
// @Summary Process checkout operation
// @Description Processes the checkout of the authenticated user's cart.
// @Tags cart
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} modelApi.CheckoutResponse "Checkout status"
// @Failure 401 {object} modelApi.ErrorResponse "Unauthorized"
// @Failure 400 {object} modelApi.ErrorResponse "Bad Request"
// @Failure 500 {object} modelApi.ErrorResponse "Internal Server Error"
// @Param input body modelApi.CheckoutRequest true "Checkout payload"
// @Router /api/cart/checkout [post]
func (h *CartHandler) Checkout(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
		return
	}

	var input modelApi.CheckoutRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	checkout, err := h.service.Checkout(c.Request.Context(), claims.UserID, model.CheckoutInput{
		AddressID: input.AddressID,
	})
	if err != nil {
		if errors.Is(err, model.ErrInvalidCart) {
			c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, converter.ToAPICheckoutFromService(checkout))
}

// CartPaymentCallback godoc
// @Summary      Process cart payment callback
// @Description  Receive payment provider callback and finalize checkout
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        payload  body      object  true  "Payment callback payload"
// @Success      200      {object}  modelApi.MessageResponse
// @Failure      400      {object}  modelApi.ErrorResponse
// @Router       /api/callback/checkout [post]
func (h *CartHandler) CartPaymentCallback(c *gin.Context) {

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "Invalid JSON"})
		return
	}
	commitCheckout, err := converter.ToServiceCheckoutFromApi(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	h.service.CommitCheckout(c.Request.Context(), commitCheckout)
	c.JSON(http.StatusOK, gin.H{"message": "JSON received"})
}

// MockCheckoutSuccess godoc
// @Summary      Mock successful checkout
// @Description  Debug endpoint that simulates successful payment for an order
// @Tags         cart
// @Produce      json
// @Param        order_id  query     string  true  "Order ID"
// @Success      302
// @Failure      400  {object}  modelApi.ErrorResponse
// @Failure      500  {object}  modelApi.ErrorResponse
// @Router       /api/cart/mock-success [get]
func (h *CartHandler) MockCheckoutSuccess(c *gin.Context) {
	orderID := c.Query("order_id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "missing order_id"})
		return
	}

	err := h.service.CommitCheckout(c.Request.Context(), model.CommitCheckout{
		OrderID:       orderID,
		PaymentStatus: model.PaymentStatusApproved,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "http://localhost:3000/orders")
}
