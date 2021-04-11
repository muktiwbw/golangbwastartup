package transaction

import (
	"bwastartup/entities/campaign"
	"bwastartup/entities/user"
	"time"
)

type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	PaymentURL string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Campaign   campaign.Campaign
	User       user.User
}
