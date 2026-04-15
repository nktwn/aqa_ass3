package handler

import (
	"diploma/modules/auth/jwt"
	"diploma/modules/user/handler/converter"
	modelApi "diploma/modules/user/handler/model"
	contextkeys "diploma/pkg/context-keys"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProfile godoc
// @Summary      Get user profile
// @Description  Fetch the authenticated user's profile
// @Tags         user
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200 {object} modelApi.GetUserProfileResponse
// @Failure      401 {object} modelApi.ErrorResponse
// @Router       /api/user/profile [get]
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok || claims == nil {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
		return
	}

	user, err := h.service.User(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return

	}

	c.JSON(http.StatusOK, modelApi.GetUserProfileResponse{
		User: converter.ToApiUserFromService(user),
	})
}

// UpdateUserProfile godoc
// @Summary      Update user profile
// @Description  Update the authenticated user's profile
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request body modelApi.UpdateUserProfileRequest true "Update Profile Request"
// @Success      200 {object} modelApi.GetUserProfileResponse
// @Failure      400 {object} modelApi.ErrorResponse
// @Failure      401 {object} modelApi.ErrorResponse
// @Failure      500 {object} modelApi.ErrorResponse
// @Router       /api/user/profile [put]
func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok || claims == nil {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
		return
	}

	var req modelApi.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	user, err := h.service.UpdateUser(c.Request.Context(), converter.ToUserFromApi(claims.UserID, req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, modelApi.GetUserProfileResponse{
		User: converter.ToApiUserFromService(user),
	})
}

// GetUserRole godoc
// @Summary      Get user role
// @Description  Returns role ID for authenticated user
// @Tags         user
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200 {object} modelApi.UserRoleResponse
// @Failure      401 {object} modelApi.ErrorResponse
// @Failure      500 {object} modelApi.ErrorResponse
// @Router       /api/user/role [get]
func (h *UserHandler) GetUserRole(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
		return
	}

	role, err := h.service.GetUserRole(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}
