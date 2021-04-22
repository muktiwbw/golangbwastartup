package handlers

import (
	"bwastartup/auth"
	"bwastartup/entities/user"
	"bwastartup/helpers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

// New UserHandler => Instanciate new UserHandler object
func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
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

	// Generate access token
	accessToken, err := h.authService.GenerateToken(newUser.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	data := helpers.APIResponse("User telah didaftarkan", 201, "created", user.FormatUser(newUser, accessToken))

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

	token, err := h.authService.GenerateToken(authenticatedUser.ID)

	if err != nil {
		data := helpers.APIResponse("Terdapat kesalahan membuat token", http.StatusInternalServerError, "error", gin.H{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, data)

		return
	}

	c.JSON(http.StatusOK, helpers.APIResponse("Berhasil login", http.StatusOK, "success", user.FormatUser(authenticatedUser, token)))
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
		c.JSON(http.StatusOK, helpers.APIResponse("Email dapat digunakan.", http.StatusOK, "ok", gin.H{"is_available": true}))

		return
	}

	c.JSON(http.StatusOK, helpers.APIResponse("Email telah digunakan.", http.StatusOK, "fail", gin.H{"is_available": false}))
}

func (h *userHandler) UpdateAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.APIResponse("Failed to extract file input", http.StatusBadRequest, "error", gin.H{"error": err.Error()}))

		return
	}

	authUser := c.MustGet("authUser").(user.User)

	driveFileID, err := h.userService.UpdateAvatar(authUser, file)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Failed to save file", http.StatusInternalServerError, "fail", gin.H{"is_uploaded": false}))

		return
	}

	c.JSON(http.StatusOK, helpers.APIResponse("Sukses menyimpan file", http.StatusOK, "ok", gin.H{"is_uploaded": true, "filename": driveFileID}))

}

func (h *userHandler) FetchCurrentUser(c *gin.Context) {
	authUser := c.MustGet("authUser").(user.User)

	c.JSON(http.StatusOK, helpers.APIResponse("Successfully fetch user data", http.StatusOK, "success", user.FormatUser(authUser, "")))
}
