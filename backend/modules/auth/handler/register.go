package handler

import (
	"log"
	"net/http"

	"diploma/pkg/validator"

	"diploma/modules/auth/handler/converter"
	modelApi "diploma/modules/auth/handler/model"
	"diploma/modules/auth/model"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary      User registration
// @Description  Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body modelApi.RegisterInput true "Register input"
// @Success      201  {object}  modelApi.RegisterResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Router       /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {

	var form struct {
		modelApi.RegisterInput
		validator.Validator
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "Name cannot be blank")
	form.CheckField(validator.NotBlank(form.PhoneNumber), "phone_number", "Phone number cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be blank")
	form.CheckField(form.Password == form.ConfirmPassword, "confirm_password", "Passwords do not match")

	if !form.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"errors": form.Errors})
		return
	}
	authUser := converter.ToServiceFromRegisterInput(form.RegisterInput)
	authUser.Info.Role = model.CustomerRole
	id, err := h.service.Register(c.Request.Context(), authUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Registered user with id: %d", id)

	c.JSON(http.StatusCreated, modelApi.RegisterResponse{
		ID: id,
	})
}
