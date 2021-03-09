package campaign

import "gorm.io/gorm"

type Repository interface {
	All() ([]Campaign, error)
	AllByUserID(userID int) ([]Campaign, error)
	Get(id int) (Campaign, error)
	// Save(campaign Campaign) (Campaign, error)
	// Update(campaign Campaign) (Campaign, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) All() ([]Campaign, error) {
	var campaigns []Campaign

	err := r.db.Find(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *repository) AllByUserID(userID int) ([]Campaign, error) {
	var campaigns []Campaign

	err := r.db.Where("user_id = ?", userID).Find(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *repository) Get(id int) (Campaign, error) {
	var campaign Campaign

	err := r.db.Where("id = ?", id).Find(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

// func (r *repository) Save(campaign Campaign) (Campaign, error) {

// }

// func (r *repository) Update(campaign Campaign) (Campaign, error) {

// }
