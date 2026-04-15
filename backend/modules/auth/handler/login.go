package handler

import (
	"net/http"

	"diploma/pkg/validator"

	modelApi "diploma/modules/auth/handler/model"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary      User login
// @Description  Login user and return tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body modelApi.LoginInput true "Login input"
// @Success      200  {object}  modelApi.LoginResponse
// @Failure      401  {object}  modelApi.ErrorResponse
// @Router       /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var form struct {
		modelApi.LoginInput
		validator.Validator
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form.CheckField(validator.NotBlank(form.PhoneNumber), "phone_number", "Phone number cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be blank")

	if !form.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"errors": form.Errors})
		return
	}

	accessToken, refreshToken, err := h.service.Login(c.Request.Context(), form.PhoneNumber, form.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, modelApi.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
