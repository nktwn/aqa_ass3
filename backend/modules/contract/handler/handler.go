package handler

import (
	"context"
	"net/http"
	"strconv"

	"diploma/modules/auth/jwt"
	contractConverter "diploma/modules/contract/handler/converter"
	apiModel "diploma/modules/contract/handler/model"
	serviceModel "diploma/modules/contract/model"
	contextkeys "diploma/pkg/context-keys"
	"github.com/gin-gonic/gin"
)

type Service interface {
	SignContract(ctx context.Context, contractID int64, role int, signature string) error
	GetContract(ctx context.Context, id int64) (*serviceModel.Contract, error)
	GetContractsByUser(ctx context.Context, userID int64) ([]*serviceModel.Contract, error) // 🔹 Добавьте эту строку
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Sign godoc
// @Summary Sign the contract
// @Description Signing the contract (by client or supplier)
// @Tags contracts
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body apiModel.SignRequest true "Contract ID and Signature"
// @Success 200 {object} map[string]string "Signature saved"
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Server error"
// @Router /api/contract/sign [post]
func (h *Handler) Sign(c *gin.Context) {
	claims := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)

	var req apiModel.SignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.SignContract(c.Request.Context(), req.ContractID, claims.Role, req.Signature)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "signature recorded"})
}

// Get godoc
// @Summary Get contract
// @Description Returns the contract by ID
// @Tags contracts
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Contract ID"
// @Success 200 {object} apiModel.ContractResponse "Contract"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Server error"
// @Router /api/contract/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	contractID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contract ID"})
		return
	}

	contract, err := h.service.GetContract(c.Request.Context(), contractID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contractConverter.ToAPI(contract))
}

// GetContractList godoc
// @Summary List of contracts for the current user
// @Tags contracts
// @Security ApiKeyAuth
// @Success 200 {array} apiModel.ContractResponse
// @Router /api/contract [get]
func (h *Handler) GetList(c *gin.Context) {
	claims := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)

	contracts, err := h.service.GetContractsByUser(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var res []apiModel.ContractResponse
	for _, c := range contracts {
		res = append(res, *contractConverter.ToAPI(c))
	}

	c.JSON(http.StatusOK, res)
}
