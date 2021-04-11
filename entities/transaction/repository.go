package transaction

import (
	"gorm.io/gorm"
)

type Repository interface {
	Save(transaction Transaction) (Transaction, error)
	Update(transaction Transaction) (Transaction, error)
	All() ([]Transaction, error)
	AllByRef(id int, field string) ([]Transaction, error)
	Get(id int) (Transaction, error)
	Verify(transaction Transaction) (Transaction, error)
	CalculateCampaignStats(campaignID int) (currentAmount int, backerCount int64, err error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Save(transaction Transaction) (Transaction, error) {
	err := r.db.Create(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) Update(transaction Transaction) (Transaction, error) {
	err := r.db.Save(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) All() ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.Preload("Campaign").Preload("User").Find(&transactions).Error

	if err != nil {
		return transactions, err
	}

	return transactions, nil

}

func (r *repository) AllByRef(id int, field string) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.Where(field+" = ?", id).Preload("Campaign").Preload("User").Find(&transactions).Error

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *repository) Get(id int) (Transaction, error) {
	var transaction Transaction

	err := r.db.Where("id = ?", id).Preload("Campaign").Preload("User").Find(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) Verify(transaction Transaction) (Transaction, error) {
	err := r.db.Save(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) CalculateCampaignStats(campaignID int) (currentAmount int, backerCount int64, err error) {
	err = r.db.Model(&Transaction{}).Where("campaign_id = ?", campaignID).Count(&backerCount).Error

	if err != nil {
		return currentAmount, backerCount, err
	}

	var trx Transaction

	err = r.db.Model(&Transaction{}).Where("campaign_id = ?", campaignID).Select("sum(amount) as amount").Scan(&trx).Error

	if err != nil {
		return currentAmount, backerCount, err
	}

	currentAmount = trx.Amount

	// fmt.Printf("Current amount: %d\nBacker count: %d\n", currentAmount, backerCount)

	return currentAmount, backerCount, nil

}
