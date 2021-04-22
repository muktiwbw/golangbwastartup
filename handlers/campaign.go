package handlers

import (
	"bwastartup/entities/campaign"
	"bwastartup/entities/user"
	"bwastartup/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
}

func NewCampaignHandler(campaignService campaign.Service) *campaignHandler {
	return &campaignHandler{campaignService}
}

func (h campaignHandler) GetAllCampaigns(c *gin.Context) {
	campaigns, err := h.campaignService.GetAllCampaigns()

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	formattedCampaigns := []campaign.CampaignThumbnailFormat{}
	// Use above instead of => var formattedCampaigns []campaign.CampaignFormat
	// Because the latter's default value would be null instead of empty slice

	for _, cmp := range campaigns {
		formattedCampaigns = append(formattedCampaigns, campaign.FormatCampaignThumbnail(cmp))
	}

	c.JSON(http.StatusOK, helpers.APIResponse("Ok", http.StatusOK, "success", formattedCampaigns))
}

func (h campaignHandler) GetOwnCampaigns(c *gin.Context) {
	// Make sure that the type is User
	authUser := c.MustGet("authUser").(user.User)

	campaigns, err := h.campaignService.GetCampaigsByUserID(authUser.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	formattedCampaigns := []campaign.CampaignThumbnailFormat{}

	for _, cmp := range campaigns {
		formattedCampaigns = append(formattedCampaigns, campaign.FormatCampaignThumbnail(cmp))
	}

	c.JSON(http.StatusOK, helpers.APIResponse("Ok", http.StatusOK, "success", formattedCampaigns))
}

func (h campaignHandler) GetCampaignByID(c *gin.Context) {
	var input campaign.GetCampaignByIDInput

	err := c.ShouldBindUri(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.APIResponse("Kesalahan pada input campaign ID", http.StatusBadRequest, "error", helpers.GetValidationErrors(err)))

		return
	}

	foundCampaign, err := h.campaignService.GetCampaignByID(input.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	if foundCampaign.ID <= 0 {
		c.JSON(http.StatusNotFound, helpers.APIResponse("Not found", http.StatusNotFound, "success", nil))

		return
	}

	c.JSON(http.StatusOK, helpers.APIResponse("Ok", http.StatusOK, "success", campaign.FormatCampaign(foundCampaign)))
}

func (h campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.APIResponse("Invalid input", http.StatusBadRequest, "error", helpers.GetValidationErrors(err)))

		return
	}

	authUser := c.MustGet("authUser").(user.User)
	input.UserID = authUser.ID

	createdCampaign, err := h.campaignService.CreateCampaign(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("There's an error from the server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	c.JSON(http.StatusCreated, helpers.APIResponse("Successfully created a campaign", http.StatusCreated, "created", campaign.FormatCampaign(createdCampaign)))
}

func (h campaignHandler) UpdateCampaign(c *gin.Context) {
	var input campaign.UpdateCampaignInput
	var uri campaign.GetCampaignByIDInput

	err := c.ShouldBindUri(&uri)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.APIResponse("Kesalahan pada input campaign ID", http.StatusBadRequest, "error", helpers.GetValidationErrors(err)))

		return
	}

	foundCampaign, err := h.campaignService.GetCampaignByID(uri.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	if foundCampaign.ID <= 0 {
		c.JSON(http.StatusNotFound, helpers.APIResponse("Data campaign tidak ditemukan", http.StatusNotFound, "not-found", nil))

		return
	}

	authUser := c.MustGet("authUser").(user.User)

	if authUser.ID != foundCampaign.UserID {
		c.JSON(http.StatusUnauthorized, helpers.APIResponse("Anda tidak punya wewenang untuk mengubah data campaign ini", http.StatusUnauthorized, "unauthorized", nil))

		return
	}

	err = c.ShouldBindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.APIResponse("Kesalahan pada input field", http.StatusBadRequest, "error", helpers.GetValidationErrors(err)))

		return
	}

	updatedCampaign, err := h.campaignService.UpdateCampaign(foundCampaign, input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	c.JSON(http.StatusCreated, helpers.APIResponse("Successfully updated a campaign", http.StatusCreated, "updated", campaign.FormatCampaign(updatedCampaign)))
}

func (h campaignHandler) DeleteCampaign(c *gin.Context) {
	var uri campaign.GetCampaignByIDInput

	err := c.ShouldBindUri(&uri)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.APIResponse("Kesalahan pada input campaign ID", http.StatusBadRequest, "error", helpers.GetValidationErrors(err)))

		return
	}

	foundCampaign, err := h.campaignService.GetCampaignByID(uri.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	if foundCampaign.ID <= 0 {
		c.JSON(http.StatusNotFound, helpers.APIResponse("Data campaign tidak ditemukan", http.StatusNotFound, "not-found", nil))

		return
	}

	authUser := c.MustGet("authUser").(user.User)

	if authUser.ID != foundCampaign.UserID {
		c.JSON(http.StatusUnauthorized, helpers.APIResponse("Anda tidak punya wewenang untuk mengubah data campaign ini", http.StatusUnauthorized, "unauthorized", nil))

		return
	}

	err = h.campaignService.DeleteCampaign(foundCampaign)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	c.JSON(http.StatusNoContent, helpers.APIResponse("Successfully deleted a campaign", http.StatusNoContent, "deleted", nil))
}

func (h campaignHandler) CreateCampaignImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("There's an error from the server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	images := form.File["images"]
	if len(images) <= 0 {
		c.JSON(http.StatusBadRequest, helpers.APIResponse("File not found", http.StatusBadRequest, "error", gin.H{"error": err.Error()}))

		return
	}

	var uri campaign.GetCampaignByIDInput
	if err = c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, helpers.APIResponse("Invalid input", http.StatusBadRequest, "error", helpers.GetValidationErrors(err)))

		return
	}

	foundCampaign, err := h.campaignService.GetCampaignByID(uri.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("There's an error from the server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	if foundCampaign.ID <= 0 {
		c.JSON(http.StatusNotFound, helpers.APIResponse("Campaign not found", http.StatusNotFound, "not-found", nil))

		return
	}

	coverIndex, err := strconv.Atoi(form.Value["cover_index"][0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("There's an error from the server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	createdImages, err := h.campaignService.CreateCampaignImages(foundCampaign.ID, coverIndex, images)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("There's an error from the server", http.StatusInternalServerError, "error", gin.H{"are_uploaded": false}))

		return
	}

	c.JSON(http.StatusCreated, helpers.APIResponse("Successfully uploaded campaign images", http.StatusCreated, "created", gin.H{"are_uploaded": true, "images": campaign.FormatCampaignImages(createdImages)}))
}
