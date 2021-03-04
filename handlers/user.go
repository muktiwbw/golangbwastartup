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

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	var input user.CheckEmailAvailabilityInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, helpers.APIResponse("Terdapat kesalahan input data", http.StatusUnprocessableEntity, "error", helpers.GetValidationErrors(err)))

		return
	}

	isAvailable, err := h.userService.EmailIsAvailable(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terdapat kesalahan input data", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	if isAvailable {
		c.JSON(http.StatusOK, helpers.APIResponse("Email tersedia.", http.StatusOK, "ok", nil))

		return
	}

	c.JSON(http.StatusOK, helpers.APIResponse("Email tidak tersedia.", http.StatusOK, "fail", nil))
}

func (h *userHandler) UpdateAvatar(c *gin.Context) {
	updatedUser, err := h.userService.SaveImage(c, "avatar", "users")

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Tidak dapat menyimpan file", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	c.JSON(http.StatusOK, helpers.APIResponse("Sukses menyimpan foto", http.StatusOK, "success", user.FormatUser(updatedUser, "secureAccessToken")))

}
