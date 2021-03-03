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

// New UserHandler => Instanciate new UserHandler object
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
		c.JSON(http.StatusInternalServerError, err.Error())

		return
	}

	data := helpers.APIResponse("User telah didaftarkan", 201, "created", user.FormatUser(newUser, "secureAccessToken"))

	c.JSON(http.StatusOK, data)
}

func (h *userHandler) LoginUser(c *gin.Context) {
	var input user.LoginUserInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		data := helpers.APIResponse("Terdapat kesalahan input data", http.StatusUnprocessableEntity, "error", helpers.GetValidationErrors(err))

		c.JSON(http.StatusUnprocessableEntity, data)

		return
	}

	// call LoginUser service
	authenticatedUser, err := h.userService.LoginUser(input)

	if err != nil {
		data := helpers.APIResponse("Terdapat kesalahan login", http.StatusUnauthorized, "error", gin.H{"error": err.Error()})
		c.JSON(http.StatusUnauthorized, data)

		return
	}

	c.JSON(http.StatusOK, user.FormatUser(authenticatedUser, "superSecureToken"))
}
