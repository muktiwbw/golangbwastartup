package handlers

import (
	"bwastartup/entities/campaign"
	"bwastartup/entities/user"
	"bwastartup/helpers"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"time"

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
		c.JSON(http.StatusBadRequest, helpers.APIResponse("Kesalahan pada input field", http.StatusBadRequest, "error", helpers.GetValidationErrors(err)))

		return
	}

	authUser := c.MustGet("authUser").(user.User)
	input.UserID = authUser.ID

	createdCampaign, err := h.campaignService.CreateCampaign(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

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
	var uri campaign.GetCampaignByIDInput

	err := c.ShouldBindUri(&uri)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.APIResponse("Kesalahan pada input field", http.StatusBadRequest, "error", helpers.GetValidationErrors(err)))

		return
	}

	foundCampaign, err := h.campaignService.GetCampaignByID(uri.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	if foundCampaign.ID <= 0 {
		c.JSON(http.StatusNotFound, helpers.APIResponse("Campaign tidak ditemukan", http.StatusNotFound, "not-found", nil))

		return
	}

	form, err := c.MultipartForm()

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	images := form.File["images"]
	coverIndex, err := strconv.Atoi(form.Value["cover_index"][0])
	// fmt.Println(images, coverIndex)
	// fmt.Println(c.FullPath())

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	campaignImages := []campaign.CampaignImage{}

UploadImageLoop:
	for i, image := range images {
		fileExt := filepath.Ext(image.Filename)
		fileTimestamp := time.Now().UnixNano()
		fileName := fmt.Sprintf("img-%d-%d%s", foundCampaign.ID, fileTimestamp, fileExt)
		fullDir := path.Join("images", "campaigns", fileName)

		err := c.SaveUploadedFile(image, fullDir)

		if err != nil {
			break UploadImageLoop
		}

		campaignImage := campaign.CampaignImage{
			CampaignID: foundCampaign.ID,
			Filename:   fileName,
			IsCover:    false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if coverIndex == i {
			campaignImage.IsCover = true
		}

		campaignImages = append(campaignImages, campaignImage)

	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	_, err = h.campaignService.CreateCampaignImages(campaignImages)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"error": err.Error()}))

		return
	}

	// Get the campaign with campaign images preloaded
	_, err = h.campaignService.GetCampaignByID(foundCampaign.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.APIResponse("Terjadi kesalahan pada server", http.StatusInternalServerError, "error", gin.H{"are_uploaded": false}))

		return
	}

	c.JSON(http.StatusCreated, helpers.APIResponse("Successfully uploaded campaign images", http.StatusCreated, "created", gin.H{"are_uploaded": true, "images": campaign.FormatCampaignImages(campaignImages)}))
}
