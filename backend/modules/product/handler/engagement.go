package handler

import (
	"diploma/modules/auth/jwt"
	"diploma/modules/product/handler/converter"
	modelApi "diploma/modules/product/handler/model"
	"diploma/modules/product/model"
	contextkeys "diploma/pkg/context-keys"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetFavorites godoc
// @Summary      Get favorite products
// @Description  Retrieve the authenticated user's favorite products
// @Tags         product
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  modelApi.FavoriteListResponse
// @Failure      401  {object}  modelApi.ErrorResponse
// @Failure      500  {object}  modelApi.ErrorResponse
// @Router       /api/product/favorites [get]
func (h *CatalogHandler) GetFavorites(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: "unauthorized"})
		return
	}

	products, err := h.service.FavoriteProducts(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	response := make([]modelApi.Product, 0, len(products))
	for _, product := range products {
		response = append(response, *converter.ToAPIProductFromService(product))
	}

	c.JSON(http.StatusOK, modelApi.FavoriteListResponse{ProductList: response})
}

// GetFavoriteStatus godoc
// @Summary      Get favorite status
// @Description  Check whether the product is in the authenticated user's favorites
// @Tags         product
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  modelApi.FavoriteStatusResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Failure      401  {object}  modelApi.ErrorResponse
// @Failure      500  {object}  modelApi.ErrorResponse
// @Router       /api/product/{id}/favorite [get]
func (h *CatalogHandler) GetFavoriteStatus(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: "unauthorized"})
		return
	}

	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "invalid product ID"})
		return
	}

	isFavorite, err := h.service.FavoriteStatus(c.Request.Context(), claims.UserID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, modelApi.FavoriteStatusResponse{IsFavorite: isFavorite})
}

// AddFavorite godoc
// @Summary      Add product to favorites
// @Description  Add the selected product to the authenticated user's favorites
// @Tags         product
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  modelApi.ProductMessageResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Failure      401  {object}  modelApi.ErrorResponse
// @Failure      500  {object}  modelApi.ErrorResponse
// @Router       /api/product/{id}/favorite [post]
func (h *CatalogHandler) AddFavorite(c *gin.Context) {
	h.updateFavorite(c, true)
}

// RemoveFavorite godoc
// @Summary      Remove product from favorites
// @Description  Remove the selected product from the authenticated user's favorites
// @Tags         product
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  modelApi.ProductMessageResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Failure      401  {object}  modelApi.ErrorResponse
// @Failure      500  {object}  modelApi.ErrorResponse
// @Router       /api/product/{id}/favorite [delete]
func (h *CatalogHandler) RemoveFavorite(c *gin.Context) {
	h.updateFavorite(c, false)
}

func (h *CatalogHandler) updateFavorite(c *gin.Context, shouldAdd bool) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: "unauthorized"})
		return
	}

	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "invalid product ID"})
		return
	}

	if shouldAdd {
		err = h.service.AddFavorite(c.Request.Context(), claims.UserID, productID)
	} else {
		err = h.service.RemoveFavorite(c.Request.Context(), claims.UserID, productID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// AddReview godoc
// @Summary      Add product review
// @Description  Add a review for the selected product on behalf of the authenticated user
// @Tags         product
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id      path      int                   true  "Product ID"
// @Param        input   body      modelApi.ReviewInput  true  "Review payload"
// @Success      200     {object}  modelApi.ProductMessageResponse
// @Failure      400     {object}  modelApi.ErrorResponse
// @Failure      401     {object}  modelApi.ErrorResponse
// @Failure      500     {object}  modelApi.ErrorResponse
// @Router       /api/product/{id}/reviews [post]
func (h *CatalogHandler) AddReview(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: "unauthorized"})
		return
	}

	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "invalid product ID"})
		return
	}

	var input modelApi.ReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	err = h.service.AddReview(c.Request.Context(), &model.ReviewInput{
		ProductID: productID,
		UserID:    claims.UserID,
		Rating:    input.Rating,
		Comment:   input.Comment,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review saved"})
}
