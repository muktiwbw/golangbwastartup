package campaign

import "time"

type Campaign struct {
	ID            int
	UserID        int
	Name          string
	Highlight     string
	Description   string
	GoalAmount    int
	CurrentAmount int
	Perks         string
	BackersCount  int
	Slug          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
