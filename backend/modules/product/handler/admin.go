package handler

import (
	"diploma/modules/auth/jwt"
	"diploma/modules/product/handler/converter"
	modelApi "diploma/modules/product/handler/model"
	"diploma/modules/product/model"
	userModel "diploma/modules/user/model"
	contextkeys "diploma/pkg/context-keys"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateProduct godoc
// @Summary      Create product
// @Description  Create a new product. Admin access required.
// @Tags         product
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        input  body      modelApi.AdminProductInput  true  "Product payload"
// @Success      201    {object}  modelApi.AdminProductResponse
// @Failure      400    {object}  modelApi.ErrorResponse
// @Failure      403    {object}  modelApi.ErrorResponse
// @Failure      500    {object}  modelApi.ErrorResponse
// @Router       /api/product/admin [post]
func (h *CatalogHandler) CreateProduct(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok || claims.Role != userModel.AdminRole {
		c.JSON(http.StatusForbidden, modelApi.ErrorResponse{Err: "admin access required"})
		return
	}

	var input modelApi.AdminProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	product, err := h.service.CreateProduct(c.Request.Context(), &model.AdminProductInput{
		Name:     input.Name,
		ImageUrl: input.ImageUrl,
		GTIN:     input.GTIN,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"product": converter.ToAPIProductFromService(product)})
}

// UpdateProduct godoc
// @Summary      Update product
// @Description  Update an existing product. Admin access required.
// @Tags         product
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id     path      int                         true  "Product ID"
// @Param        input  body      modelApi.AdminProductInput  true  "Product payload"
// @Success      200    {object}  modelApi.AdminProductResponse
// @Failure      400    {object}  modelApi.ErrorResponse
// @Failure      403    {object}  modelApi.ErrorResponse
// @Failure      500    {object}  modelApi.ErrorResponse
// @Router       /api/product/admin/{id} [put]
func (h *CatalogHandler) UpdateProduct(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok || claims.Role != userModel.AdminRole {
		c.JSON(http.StatusForbidden, modelApi.ErrorResponse{Err: "admin access required"})
		return
	}

	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "invalid product ID"})
		return
	}

	var input modelApi.AdminProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	product, err := h.service.UpdateProduct(c.Request.Context(), &model.AdminProductInput{
		ID:       productID,
		Name:     input.Name,
		ImageUrl: input.ImageUrl,
		GTIN:     input.GTIN,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": converter.ToAPIProductFromService(product)})
}

// DeleteProduct godoc
// @Summary      Delete product
// @Description  Delete an existing product. Admin access required.
// @Tags         product
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  modelApi.ProductMessageResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Failure      403  {object}  modelApi.ErrorResponse
// @Failure      500  {object}  modelApi.ErrorResponse
// @Router       /api/product/admin/{id} [delete]
func (h *CatalogHandler) DeleteProduct(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok || claims.Role != userModel.AdminRole {
		c.JSON(http.StatusForbidden, modelApi.ErrorResponse{Err: "admin access required"})
		return
	}

	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "invalid product ID"})
		return
	}

	if err := h.service.DeleteProduct(c.Request.Context(), productID); err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}
