package order

import (
	"diploma/modules/order/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h *handler.OrderHandler) {
	orderRoutes := router.Group("/order")
	{
		// orderRoutes.POST("", h.CreateOrder)
		orderRoutes.GET("", h.GetOrders)
		orderRoutes.GET("/:id", h.GetOrderByID)
		orderRoutes.POST("/status", h.UpdateOrderStatus)
		orderRoutes.POST("/cancel", h.CancelOrder)
	}
}
