package user

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Save(user User) (User, error)
	FindByEmail(email string) (User, error)
	FindByID(id int) (User, error)
	Update(user User) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindByEmail(email string) (User, error) {
	var foundUser User

	err := r.db.Where("email = ?", email).Find(&foundUser).Error

	if err != nil {
		return foundUser, err
	}

	if foundUser.ID == 0 {
		return foundUser, errors.New("Email atau password salah.")
	}

	// fmt.Println("PASSED CONDITION")
	return foundUser, nil
}

func (r *repository) FindByID(id int) (User, error) {
	var foundUser User

	err := r.db.Where("id = ?", id).Find(&foundUser).Error

	if err != nil {
		return foundUser, err
	}

	if foundUser.ID == 0 {
		return foundUser, errors.New("User tidak ditemukan.")
	}

	// fmt.Println("PASSED CONDITION")
	return foundUser, nil
}

func (r *repository) Update(user User) (User, error) {
	// Segala sesuatu dari db hanya menerima pointer sebagai argumen
	err := r.db.Save(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}
