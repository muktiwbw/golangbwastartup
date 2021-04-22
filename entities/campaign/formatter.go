package campaign

import (
	"time"

	"github.com/muktiwbw/gdstorage"
)

type CampaignFormat struct {
	ID            int                       `json:"id"`
	Name          string                    `json:"name"`
	Highlight     string                    `json:"highlight"`
	Description   string                    `json:"description"`
	CoverImage    string                    `json:"cover_image"`
	Images        []CampaignImageFormat     `json:"images"`
	GoalAmount    int                       `json:"goal_amount"`
	CurrentAmount int                       `json:"current_amount"`
	BackersCount  int                       `json:"backers_count"`
	Perks         string                    `json:"perks"`
	User          CampaignUserSnippetFormat `json:"user"`
	CreatedAt     time.Time                 `json:"created_at"`
}

type CampaignUserSnippetFormat struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type CampaignThumbnailFormat struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Highlight     string    `json:"highlight"`
	Image         string    `json:"image"`
	GoalAmount    int       `json:"goal_amount"`
	CurrentAmount int       `json:"current_amount"`
	BackersCount  int       `json:"backers_count"`
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
	var coverImage string
	var avatar string

ImageFilterLoop:
	for _, image := range campaign.CampaignImages {
		if image.IsCover {
			coverImage = gdstorage.GetURL(image.Filename)

			break ImageFilterLoop
		}
	}

	if campaign.User.Avatar != "" {
		avatar = gdstorage.GetURL(campaign.User.Avatar)
	}

	return CampaignFormat{
		ID:            campaign.ID,
		Name:          campaign.Name,
		Highlight:     campaign.Highlight,
		Description:   campaign.Description,
		CoverImage:    coverImage,
		Images:        FormatCampaignImages(campaign.CampaignImages),
		GoalAmount:    campaign.GoalAmount,
		CurrentAmount: campaign.CurrentAmount,
		BackersCount:  campaign.BackersCount,
		Perks:         campaign.Perks,
		User: CampaignUserSnippetFormat{
			ID:     campaign.User.ID,
			Name:   campaign.User.Name,
			Avatar: avatar,
		},
		CreatedAt: campaign.CreatedAt,
	}
}

func FormatCampaignThumbnail(campaign Campaign) CampaignThumbnailFormat {
	var image string

	if len(campaign.CampaignImages) >= 1 {

	ImageFilterLoop:
		for _, i := range campaign.CampaignImages {
			if i.IsCover {
				image = gdstorage.GetURL(i.Filename)

				break ImageFilterLoop
			}
		}
	}

	return CampaignThumbnailFormat{
		ID:            campaign.ID,
		Name:          campaign.Name,
		Highlight:     campaign.Highlight,
		Image:         image,
		GoalAmount:    campaign.GoalAmount,
		CurrentAmount: campaign.CurrentAmount,
		BackersCount:  campaign.BackersCount,
		UserID:        campaign.UserID,
		CreatedAt:     campaign.CreatedAt,
	}
}

func FormatCampaignImages(images []CampaignImage) []CampaignImageFormat {
	formattedCampaignImages := []CampaignImageFormat{}

	for _, image := range images {
		formattedCampaignImages = append(formattedCampaignImages, CampaignImageFormat{Filename: gdstorage.GetURL(image.Filename), IsCover: image.IsCover})
	}

	return formattedCampaignImages
}
