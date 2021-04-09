package campaign

import (
	"fmt"
	"time"

	"github.com/gosimple/slug"
)

type Service interface {
	GetAllCampaigns() ([]Campaign, error)
	GetCampaigsByUserID(userID int) ([]Campaign, error)
	GetCampaignByID(id int) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(campaign Campaign, updateValues UpdateCampaignInput) (Campaign, error)
	DeleteCampaign(campaign Campaign) error
	CreateCampaignImages(images []CampaignImage) ([]CampaignImage, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) GetAllCampaigns() ([]Campaign, error) {
	campaigns, err := s.repository.All()

	if err != nil {
		return campaigns, err
	}

	return campaigns, err
}

func (s *service) GetCampaigsByUserID(userID int) ([]Campaign, error) {
	campaigns, err := s.repository.AllByUserID(userID)

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (s *service) GetCampaignByID(id int) (Campaign, error) {
	campaign, err := s.repository.Get(id)

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{
		UserID:        input.UserID,
		Name:          input.Name,
		Highlight:     input.Highlight,
		Description:   input.Description,
		GoalAmount:    input.GoalAmount,
		CurrentAmount: 0,
		Perks:         input.Perks,
		BackersCount:  0,
		Slug:          slug.Make(fmt.Sprintf("%d %s", input.UserID, input.Name)),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	newCampaign, err := s.repository.Save(campaign)

	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}

func (s *service) UpdateCampaign(campaign Campaign, updateValues UpdateCampaignInput) (Campaign, error) {
	newValuesCampaign := Campaign{
		Name:          updateValues.Name,
		Highlight:     updateValues.Highlight,
		Description:   updateValues.Description,
		GoalAmount:    updateValues.GoalAmount,
		Perks:         updateValues.Perks,
		UpdatedAt:     time.Now(),
		CurrentAmount: updateValues.CurrentAmount,
		BackersCount:  int(updateValues.BackersCount),
	}

	updatedCampaign, err := s.repository.Update(campaign, newValuesCampaign)

	if err != nil {
		return updatedCampaign, err
	}

	return updatedCampaign, nil
}

func (s *service) DeleteCampaign(campaign Campaign) error {
	err := s.repository.Delete(campaign)

	if err != nil {
		return err
	}

	return nil
}

func (s *service) CreateCampaignImages(images []CampaignImage) ([]CampaignImage, error) {
	var err error = nil

	createdImages := []CampaignImage{}

	err = s.repository.ResetCampaignImageCover(images[0].CampaignID)

	if err != nil {
		return createdImages, err
	}

SavingEachImage: // It's a looping label, later useful for breaking out of a loop
	for _, image := range images {
		img, err := s.repository.SaveCampaignImage(image)

		if err != nil {
			break SavingEachImage
		}

		createdImages = append(createdImages, img)
	}

	if err != nil {
		return createdImages, err
	}

	return createdImages, nil
}
