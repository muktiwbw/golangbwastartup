package user

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	LoginUser(input LoginUserInput) (User, error)
	EmailIsAvailable(input CheckEmailAvailabilityInput) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	u := User{}

	u.Name = input.Name
	u.Email = input.Email
	u.Occupation = input.Occupation
	u.Avatar = "default.jpg"
	u.Role = "user"
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	// ===================== Hash Password =====================
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil {
		return u, err
	}

	u.Password = string(hashedPassword)
	// ===================== Hash Password =====================

	// ===================== Save User Data =====================
	nu, err := s.repository.Save(u)
	// ===================== Save User Data =====================

	if err != nil {
		return nu, errors.New("Terjadi kesalahan ketika menyimpan data.")
	}

	return nu, nil
}

func (s *service) LoginUser(input LoginUserInput) (User, error) {
	// call Login repository
	authenticatedUser, err := s.repository.FindByEmail(input.Email)

	if err != nil {
		return authenticatedUser, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(authenticatedUser.Password), []byte(input.Password))

	if err != nil {
		// Sengaja dijadikan sama dengan "email not found" supaya user tidak tau mana yang salah
		return authenticatedUser, errors.New("Email atau password salah.")
	}

	return authenticatedUser, nil
}

func (s *service) EmailIsAvailable(input CheckEmailAvailabilityInput) (bool, error) {
	foundUser, err := s.repository.FindByEmail(input.Email)

	if err != nil {
		if foundUser.ID == 0 {
			return true, nil
		}

		return false, err
	}

	return false, nil
}
