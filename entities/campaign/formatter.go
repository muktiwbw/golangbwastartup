package campaign

type CampaignFormat struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Highlight     string `json:"highlight"`
	ImageURL      string `json:"image_url"`
	GoalAmount    int    `json:"goal_amount"`
	CurrentAmount int    `json:"current_amount"`
	UserID        int    `json:"user_id"`
}

func FormatCampaign(campaign Campaign) CampaignFormat {
	return CampaignFormat{
		ID:            campaign.ID,
		Name:          campaign.Name,
		Highlight:     campaign.Highlight,
		ImageURL:      "",
		GoalAmount:    campaign.GoalAmount,
		CurrentAmount: campaign.CurrentAmount,
		UserID:        campaign.UserID,
	}
}
