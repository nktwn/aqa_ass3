package handler

import (
	"net/http"
	"strconv"
	"strings"

	"diploma/modules/product/handler/converter"
	modelApi "diploma/modules/product/handler/model"
	"diploma/modules/product/model"

	"github.com/gin-gonic/gin"
)

// GetProductList godoc
// @Summary      Get product list
// @Description  Retrieve a list of products with pagination support using limit and offset
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        limit     query     int     false "Limit number of products"  // Define limit as a query parameter
// @Param        offset    query     int     false "Offset for pagination"    // Define offset as a query parameter
// @Success      200  {object}  modelApi.ProductListResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Router       /api/product/list [get]
func (h *CatalogHandler) GetProductList(c *gin.Context) {
	// Extracting query parameters
	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0") // Default to 0 if not provided
	minPrice := c.DefaultQuery("min_price", "0")
	maxPrice := c.DefaultQuery("max_price", "0")

	// Convert limit and offset from string to int
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	minPriceInt, err := strconv.Atoi(minPrice)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min_price parameter"})
		return
	}

	maxPriceInt, err := strconv.Atoi(maxPrice)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max_price parameter"})
		return
	}

	// Prepare the input for service (query parameters can be passed as part of the input)
	input := modelApi.ProductListInput{
		Limit:     limitInt,
		Offset:    offsetInt,
		Search:    strings.TrimSpace(c.Query("search")),
		MinPrice:  minPriceInt,
		MaxPrice:  maxPriceInt,
		SortBy:    strings.TrimSpace(c.DefaultQuery("sort_by", "newest")),
		SortOrder: strings.TrimSpace(c.DefaultQuery("sort_order", "desc")),
	}

	// Call the service layer to get the product list
	productList, err := h.service.ProductList(c.Request.Context(), converter.ToServiceProductListQueryFromAPI(&input))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert service response to API response and return it
	c.JSON(http.StatusOK, converter.ToProductListResponeFromService(productList))
}

// GetFeaturedProducts godoc
// @Summary      Get featured products
// @Description  Retrieve a curated list of featured products
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        limit  query     int  false "Limit number of featured products"
// @Success      200    {object}  modelApi.ProductListResponse
// @Failure      400    {object}  modelApi.ErrorResponse
// @Failure      500    {object}  modelApi.ErrorResponse
// @Router       /api/product/featured [get]
func (h *CatalogHandler) GetFeaturedProducts(c *gin.Context) {
	limit := c.DefaultQuery("limit", "8")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	productList, err := h.service.FeaturedProducts(c.Request.Context(), limitInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, converter.ToProductListResponeFromService(productList))
}

// GetSuggestions godoc
// @Summary      Get product suggestions
// @Description  Retrieve product search suggestions by query string
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        q      query     string  false "Search query"
// @Param        limit  query     int     false "Limit number of suggestions"
// @Success      200    {object}  modelApi.ProductSuggestionResponse
// @Failure      400    {object}  modelApi.ErrorResponse
// @Failure      500    {object}  modelApi.ErrorResponse
// @Router       /api/product/suggest [get]
func (h *CatalogHandler) GetSuggestions(c *gin.Context) {
	limit := c.DefaultQuery("limit", "5")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	suggestions, err := h.service.Suggestions(c.Request.Context(), &model.SuggestionQuery{
		Query: strings.TrimSpace(c.Query("q")),
		Limit: limitInt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, modelApi.ProductSuggestionResponse{Suggestions: suggestions})
}
