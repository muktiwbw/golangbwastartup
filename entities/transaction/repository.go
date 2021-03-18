package transaction

import "gorm.io/gorm"

type Repository interface {
	Save(transaction Transaction) (Transaction, error)
	All() ([]Transaction, error)
	AllByRef(id int, field string) ([]Transaction, error)
	Get(id int) (Transaction, error)
	Verify(transaction Transaction) (Transaction, error)
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

	returnedTransaction, err := r.Get(transaction.ID)

	if err != nil {
		return returnedTransaction, err
	}

	return returnedTransaction, nil
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
