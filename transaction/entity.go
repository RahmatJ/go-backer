package transaction

import (
	"go-backer/campaign"
	"go-backer/user"
	"time"
)

type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	Campaign   campaign.Campaign
	CreatedAt  time.Time
	UpdatedAt  time.Time
	User       user.User
}
