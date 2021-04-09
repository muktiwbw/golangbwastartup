package campaign

import (
	"bwastartup/entities/user"
	"time"
)

type Campaign struct {
	ID             int
	UserID         int
	Name           string
	Highlight      string
	Description    string
	GoalAmount     int
	CurrentAmount  int
	Perks          string
	BackersCount   int
	Slug           string
	CampaignImages []CampaignImage
	User           user.User
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CampaignImage struct {
	ID         int
	CampaignID int
	Filename   string
	IsCover    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
