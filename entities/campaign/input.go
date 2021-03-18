package campaign

type GetCampaignByIDInput struct {
	ID int `uri:"campaign_id" binding:"required"`
}

type CreateCampaignInput struct {
	UserID      int
	Name        string `json:"name" binding:"required"`
	Highlight   string `json:"highlight" binding:"required"`
	Description string `json:"description" binding:"required"`
	GoalAmount  int    `json:"goal_amount" binding:"required"`
	Perks       string `json:"perks" binding:"required"`
}

type UpdateCampaignInput struct {
	ID          int
	Name        string `json:"name"`
	Highlight   string `json:"highlight"`
	Description string `json:"description"`
	GoalAmount  int    `json:"goal_amount"`
	Perks       string `json:"perks"`
}
