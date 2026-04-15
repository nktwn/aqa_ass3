package handler

import (
	"context"
	"diploma/modules/auth/jwt"
	notificationModel "diploma/modules/notification/model"
	contextkeys "diploma/pkg/context-keys"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Service interface {
	ListByUser(ctx context.Context, userID int64) ([]notificationModel.Notification, error)
	UnreadCount(ctx context.Context, userID int64) (int, error)
	MarkRead(ctx context.Context, userID, notificationID int64) error
	MarkAllRead(ctx context.Context, userID int64) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// List godoc
// @Summary      Get notifications
// @Description  Retrieve notifications for the authenticated user
// @Tags         notifications
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  ListResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/notifications [get]
func (h *Handler) List(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	notifications, err := h.service.ListByUser(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]Notification, 0, len(notifications))
	for _, notification := range notifications {
		response = append(response, Notification{
			ID:        notification.ID,
			Title:     notification.Title,
			Message:   notification.Message,
			Kind:      notification.Kind,
			IsRead:    notification.IsRead,
			OrderID:   notification.OrderID,
			CreatedAt: notification.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	c.JSON(http.StatusOK, ListResponse{Notifications: response})
}

// UnreadCount godoc
// @Summary      Get unread notifications count
// @Description  Retrieve unread notifications count for the authenticated user
// @Tags         notifications
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  UnreadCountResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/notifications/unread-count [get]
func (h *Handler) UnreadCount(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	count, err := h.service.UnreadCount(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UnreadCountResponse{Count: count})
}

// MarkRead godoc
// @Summary      Mark notification as read
// @Description  Mark a single notification as read for the authenticated user
// @Tags         notifications
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Notification ID"
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/notifications/{id}/read [post]
func (h *Handler) MarkRead(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	notificationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
		return
	}

	if err := h.service.MarkRead(c.Request.Context(), claims.UserID, notificationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification marked as read"})
}

// MarkAllRead godoc
// @Summary      Mark all notifications as read
// @Description  Mark all notifications as read for the authenticated user
// @Tags         notifications
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  MessageResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/notifications/read-all [post]
func (h *Handler) MarkAllRead(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.service.MarkAllRead(c.Request.Context(), claims.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all notifications marked as read"})
}
