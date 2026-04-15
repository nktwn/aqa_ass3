package handler

import (
	"diploma/modules/auth/jwt"
	"diploma/modules/cart/handler/converter"
	modelApi "diploma/modules/cart/handler/model"
	"diploma/modules/cart/model"
	contextkeys "diploma/pkg/context-keys"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary      Put product to Card
// @Description  --
// @Tags         cart
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        input body modelApi.AddProductToCartInput true "Put Card input"
// @Success      200  {object}  modelApi.AddProductToCardResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Router       /api/cart/add [post]
func (h *CartHandler) AddProductToCard(c *gin.Context) {

	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)

	if !ok || claims == nil {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
		return
	}

	var input modelApi.AddProductToCartInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err.Error()})
		return
	}
	input.CustomerID = claims.UserID

	fmt.Println(input)
	err := h.service.AddProductToCard(c.Request.Context(), converter.ToServiceCardInputFromAPI(&input))
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, modelApi.AddProductToCardResponse{Status: "ok"})
}

// @Summary      get cart
// @Description  --
// @Tags         cart
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object} modelApi.GetCartResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Router       /api/cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
		return
	}

	cart, err := h.service.Cart(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}
	c.JSON(http.StatusOK, converter.ToAPIGetCartFromService(cart))
}

// DeleteProductFromCart godoc
// @Summary      Delete product from cart
// @Description  Deletes given quantity of product by product_id and supplier_id from cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        product_id   query     int  true  "Product ID"
// @Param        supplier_id  query     int  true  "Supplier ID"
// @Param        quantity     query     int  false "Quantity to delete (default 1)"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/cart/delete [delete]
func (h *CartHandler) DeleteProductFromCart(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	productIDStr := c.Query("product_id")
	supplierIDStr := c.Query("supplier_id")
	quantityStr := c.DefaultQuery("quantity", "1")

	if productIDStr == "" || supplierIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing product_id or supplier_id"})
		return
	}

	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	supplierID, err := strconv.ParseInt(supplierIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplier_id"})
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quantity"})
		return
	}

	query := &model.PutCartQuery{
		CustomerID: claims.UserID,
		ProductID:  productID,
		SupplierID: supplierID,
		Quantity:   quantity,
	}

	err = h.service.DeleteProductFromCart(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted from cart"})
}



// ClearCart godoc
// @Summary      Clear cart
// @Description  Deletes all products from the user's cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  modelApi.AddProductToCardResponse
// @Failure      401  {object}  modelApi.ErrorResponse
// @Failure      500  {object}  modelApi.ErrorResponse
// @Router       /api/cart/clear [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
		return
	}

	err := h.service.ClearCart(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, modelApi.AddProductToCardResponse{Status: "ok"})
}