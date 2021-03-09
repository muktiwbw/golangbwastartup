package handlers

import (
	"bwastartup/entities/campaign"
	"bwastartup/entities/user"
	"bwastartup/helpers"
	"net/http"

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

	formattedCampaigns := []campaign.CampaignFormat{}
	// Use above instead of => var formattedCampaigns []campaign.CampaignFormat
	// Because the latter's default value would be null instead of empty slice

	for _, cmp := range campaigns {
		formattedCampaigns = append(formattedCampaigns, campaign.FormatCampaign(cmp))
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

	formattedCampaigns := []campaign.CampaignFormat{}

	for _, cmp := range campaigns {
		formattedCampaigns = append(formattedCampaigns, campaign.FormatCampaign(cmp))
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
