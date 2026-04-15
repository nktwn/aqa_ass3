package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"diploma/modules/product/handler/converter"
	modelApi "diploma/modules/product/handler/model"

	"github.com/gin-gonic/gin"
)

// GetProduct godoc
// @Summary      Get product by ID
// @Description  Retrieve product information by its ID
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id     path      int     true "Product ID"
// @Success      200  {object}  modelApi.ProductResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Failure      404  {object}  modelApi.ErrorResponse
// @Router       /api/product/{id} [get]
func (h *CatalogHandler) GetProduct(c *gin.Context) {

	productID := c.Param("id")

	// Convert limit and offset from string to int
	productIdInt, err := strconv.ParseInt(productID, 10, 64)
	fmt.Println(productIdInt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	input := modelApi.ProductInput{ID: productIdInt}
	product, err := h.service.Product(c.Request.Context(), converter.ToServieProductQueryFromApi(input))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, converter.ToApiProductResponeFromService(product))
}
