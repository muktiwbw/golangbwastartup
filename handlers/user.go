package handlers

import (
	"bwastartup/entities/user"
	"bwastartup/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput

	// Input has kilobytes of data, which makes it more sense to pass a pointer
	err := c.ShouldBindJSON(&input)

	// Check for input validation
	if err != nil {
		data := helpers.APIResponse("Terdapat kesalahan pada input", http.StatusUnprocessableEntity, "error", helpers.GetValidationErrors(err))
		c.JSON(http.StatusBadRequest, data)

		return
	}

	newUser, err := h.userService.RegisterUser(input)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	data := helpers.APIResponse("User telah didaftarkan", 201, "created", user.FormatUser(newUser, "secureAccessToken"))

	c.JSON(http.StatusOK, data)
}
