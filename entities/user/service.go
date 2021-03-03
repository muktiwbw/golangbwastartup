package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
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
		return nu, err
	}

	return nu, nil
}
