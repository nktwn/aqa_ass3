package product

import (
	"diploma/modules/product/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h *handler.CatalogHandler) {
	catalogRoutes := router.Group("product")
	{
		catalogRoutes.GET("/list", h.GetProductList)
		catalogRoutes.GET("/featured", h.GetFeaturedProducts)
		catalogRoutes.GET("/suggest", h.GetSuggestions)

		catalogRoutes.GET("/:id", h.GetProduct)

		// catalogRoutes.GET("/product/pages", h.GetPageCount)

		// catalogRoutes.POST("/product", h.AddProduct)
	}

}

func RegisterProtectedRoutes(router *gin.RouterGroup, h *handler.CatalogHandler) {
	catalogRoutes := router.Group("product")
	{
		catalogRoutes.GET("/favorites", h.GetFavorites)
		catalogRoutes.GET("/:id/favorite", h.GetFavoriteStatus)
		catalogRoutes.POST("/:id/favorite", h.AddFavorite)
		catalogRoutes.DELETE("/:id/favorite", h.RemoveFavorite)
		catalogRoutes.POST("/:id/reviews", h.AddReview)
	}
}

func RegisterAdminRoutes(router *gin.RouterGroup, h *handler.CatalogHandler) {
	adminRoutes := router.Group("product/admin")
	{
		adminRoutes.POST("", h.CreateProduct)
		adminRoutes.PUT("/:id", h.UpdateProduct)
		adminRoutes.DELETE("/:id", h.DeleteProduct)
	}
}
