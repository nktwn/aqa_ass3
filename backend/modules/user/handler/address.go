package handler

import (
	"context"
	"diploma/modules/auth/jwt"
	"diploma/modules/user/handler/converter"
	modelApi "diploma/modules/user/handler/model"
	"diploma/modules/user/model"
	contextkeys "diploma/pkg/context-keys"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IAddressService interface {
	SetAddress(ctx context.Context, input model.Address) (int64, error)
	Address(ctx context.Context, userID int64) ([]model.Address, error)
}

// Register godoc
// @Summary      Set address
// @Description  --
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        input body modelApi.SetAddressInput true "something"
// @Success      200  {object}  modelApi.StatusResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Failure      401  {object}  modelApi.ErrorResponse
// @Failure      500  {object}  modelApi.ErrorResponse
// @Router       /api/user/address [post]
func (h *UserHandler) SetAddress(c *gin.Context) {

	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)

	if !ok || claims == nil {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
		return
	}

	var input modelApi.SetAddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err.Error()})
		return
	}
	_, err := h.service.SetAddress(c.Request.Context(), converter.ToServiceAddressFromApi(claims.UserID, input.Address))
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetAddress godoc
// @Summary      Get address
// @Description  Retrieve address list for a user
// @Tags         user
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200 {object} modelApi.GetAddressResponse
// @Failure      401 {object} modelApi.ErrorResponse
// @Failure      500 {object} modelApi.ErrorResponse
// @Router       /api/user/address [get]
func (h *UserHandler) GetAddress(c *gin.Context) {

	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)

	if !ok || claims == nil {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
		return
	}

	addressList, err := h.service.Address(c.Request.Context(), claims.UserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, modelApi.GetAddressResponse{AddressList: converter.ToApiAddressListFromService(addressList)})
}
