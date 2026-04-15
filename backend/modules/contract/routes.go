package contract

import (
	"diploma/modules/contract/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h *handler.Handler) {
	routes := router.Group("/contract")
	{
		routes.POST("/sign", h.Sign)
		routes.GET("/:id", h.Get)
		routes.GET("", h.GetList) // /api/contract

	}
}
