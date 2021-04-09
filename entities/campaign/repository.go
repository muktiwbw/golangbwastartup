package campaign

import (
	"gorm.io/gorm"
)

type Repository interface {
	All() ([]Campaign, error)
	AllByUserID(userID int) ([]Campaign, error)
	Get(campaignID int) (Campaign, error)
	Save(campaign Campaign) (Campaign, error)
	Update(campaign Campaign, newValuesCampaign Campaign) (Campaign, error)
	Delete(campaign Campaign) error
	ResetCampaignImageCover(campaignID int) error
	SaveCampaignImage(image CampaignImage) (CampaignImage, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) All() ([]Campaign, error) {
	var campaigns []Campaign

	err := r.db.Preload("CampaignImages", "is_cover = true").Find(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *repository) AllByUserID(userID int) ([]Campaign, error) {
	var campaigns []Campaign

	err := r.db.Preload("CampaignImages", "is_cover = true").Where("user_id = ?", userID).Find(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *repository) Get(campaignID int) (Campaign, error) {
	var campaign Campaign

	err := r.db.Where("id = ?", campaignID).Preload("CampaignImages").Preload("User").Find(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) Save(campaign Campaign) (Campaign, error) {
	err := r.db.Create(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) Update(campaign Campaign, newValuesCampaign Campaign) (Campaign, error) {
	err := r.db.Model(&campaign).Updates(newValuesCampaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) Delete(campaign Campaign) error {
	err := r.db.Delete(&campaign).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) ResetCampaignImageCover(campaignID int) error {
	err := r.db.Model(&CampaignImage{}).Where("campaign_id = ?", campaignID).Update("is_cover", false).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) SaveCampaignImage(image CampaignImage) (CampaignImage, error) {
	err := r.db.Create(&image).Error

	if err != nil {
		return image, err
	}

	return image, nil
}
