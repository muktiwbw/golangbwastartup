package campaign

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gosimple/slug"
	"github.com/muktiwbw/gdstorage"
)

type Service interface {
	GetAllCampaigns() ([]Campaign, error)
	GetCampaigsByUserID(userID int) ([]Campaign, error)
	GetCampaignByID(id int) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(campaign Campaign, updateValues UpdateCampaignInput) (Campaign, error)
	DeleteCampaign(campaign Campaign) error
	CreateCampaignImages(campaignID int, coverIndex int, files []*multipart.FileHeader) ([]CampaignImage, error)
}

type service struct {
	repository Repository
	gds        gdstorage.GoogleDriveStorage
}

func NewService(repository Repository, gds gdstorage.GoogleDriveStorage) Service {
	return &service{repository, gds}
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

func (s *service) CreateCampaignImages(campaignID int, coverIndex int, files []*multipart.FileHeader) ([]CampaignImage, error) {
	// * Store images to Google Ddrive
	driveFileInputs := []*gdstorage.StoreFileInput{}

	for _, file := range files {
		fileExt := filepath.Ext(file.Filename)
		fileName := fmt.Sprintf("campaign-%d%s", campaignID, fileExt)

		driveFileInputs = append(driveFileInputs, &gdstorage.StoreFileInput{Name: fileName, FileHeader: file})
	}

	driveFileIDs, err := s.gds.StoreFiles(driveFileInputs, os.Getenv("DRIVE_APP_CAMPAIGN_IMAGES_DIR_ID"))
	if err != nil {
		return []CampaignImage{}, err
	}

	// * Save images data to DB
	campaignImages := []CampaignImage{}

	for i, driveFileID := range driveFileIDs {
		campaignImage := CampaignImage{}
		campaignImage.CampaignID = campaignID
		campaignImage.Filename = driveFileID
		campaignImage.IsCover = false
		campaignImage.CreatedAt = time.Now()
		campaignImage.UpdatedAt = time.Now()

		if i == coverIndex {
			campaignImage.IsCover = true
		}

		campaignImages = append(campaignImages, campaignImage)
	}

	if err := s.repository.ResetCampaignImageCover(campaignID); err != nil {
		return []CampaignImage{}, err
	}

	createdImages, err := s.repository.SaveCampaignImages(campaignImages)
	if err != nil {
		return campaignImages, err
	}

	return createdImages, nil
}
