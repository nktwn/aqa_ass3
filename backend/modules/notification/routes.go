package notification

import (
	"diploma/modules/notification/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h *handler.Handler) {
	notificationRoutes := router.Group("/notifications")
	{
		notificationRoutes.GET("", h.List)
		notificationRoutes.GET("/unread-count", h.UnreadCount)
		notificationRoutes.POST("/:id/read", h.MarkRead)
		notificationRoutes.POST("/read-all", h.MarkAllRead)
	}
}
