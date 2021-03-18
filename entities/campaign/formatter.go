package campaign

import "time"

type CampaignFormat struct {
	ID            int                   `json:"id"`
	Name          string                `json:"name"`
	Highlight     string                `json:"highlight"`
	Description   string                `json:"description"`
	Images        []CampaignImageFormat `json:"images"`
	GoalAmount    int                   `json:"goal_amount"`
	CurrentAmount int                   `json:"current_amount"`
	Perks         string                `json:"perks"`
	UserID        int                   `json:"user_id"`
	CreatedAt     time.Time             `json:"created_at"`
}

type CampaignThumbnailFormat struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Highlight     string    `json:"highlight"`
	Image         string    `json:"image"`
	GoalAmount    int       `json:"goal_amount"`
	CurrentAmount int       `json:"current_amount"`
	UserID        int       `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type CampaignSnippetFormat struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Highlight string `json:"highlight"`
}

type CampaignImageFormat struct {
	Filename string `json:"filename"`
	IsCover  bool   `json:"is_cover"`
}

func FormatCampaign(campaign Campaign) CampaignFormat {
	return CampaignFormat{
		ID:            campaign.ID,
		Name:          campaign.Name,
		Highlight:     campaign.Highlight,
		Description:   campaign.Description,
		Images:        FormatCampaignImages(campaign.CampaignImages),
		GoalAmount:    campaign.GoalAmount,
		CurrentAmount: campaign.CurrentAmount,
		Perks:         campaign.Perks,
		UserID:        campaign.UserID,
		CreatedAt:     campaign.CreatedAt,
	}
}

func FormatCampaignThumbnail(campaign Campaign) CampaignThumbnailFormat {
	var campaignImage string

	if len(campaign.CampaignImages) >= 1 {
		campaignImage = campaign.CampaignImages[0].Filename
	}

	return CampaignThumbnailFormat{
		ID:            campaign.ID,
		Name:          campaign.Name,
		Highlight:     campaign.Highlight,
		Image:         campaignImage,
		GoalAmount:    campaign.GoalAmount,
		CurrentAmount: campaign.CurrentAmount,
		UserID:        campaign.UserID,
		CreatedAt:     campaign.CreatedAt,
	}
}

func FormatCampaignImages(images []CampaignImage) []CampaignImageFormat {
	formattedCampaignImages := []CampaignImageFormat{}

	for _, image := range images {
		formattedCampaignImages = append(formattedCampaignImages, CampaignImageFormat{Filename: image.Filename, IsCover: image.IsCover})
	}

	return formattedCampaignImages
}
