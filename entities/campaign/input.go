package campaign

type GetCampaignByIDInput struct {
	ID int `uri:"campaign_id" binding:"required"`
}
