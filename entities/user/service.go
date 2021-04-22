package user

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/muktiwbw/gdstorage"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetUserByID(user_id int) (User, error)
	RegisterUser(input RegisterUserInput) (User, error)
	LoginUser(input LoginUserInput) (User, error)
	EmailIsAvailable(input CheckEmailAvailabilityInput) (bool, error)
	UpdateAvatar(user User, file *multipart.FileHeader) (string, error)
}

type service struct {
	repository Repository
	gds        gdstorage.GoogleDriveStorage
}

func NewService(repository Repository, gds gdstorage.GoogleDriveStorage) Service {
	return &service{repository, gds}
}

func (s *service) GetUserByID(user_id int) (User, error) {
	foundUser, err := s.repository.FindByID(user_id)

	if err != nil {
		return foundUser, err
	}

	return foundUser, nil
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	u := User{}

	u.Name = input.Name
	u.Email = input.Email
	u.Occupation = input.Occupation
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

func (s *service) UpdateAvatar(user User, file *multipart.FileHeader) (string, error) {
	// * Check if avatar exists
	if user.Avatar != "" {
		if err := s.gds.DeleteFile(user.Avatar); err != nil {
			if (strings.Split(err.Error(), ": "))[0] != "ERR404" {
				return "", errors.New(fmt.Sprintf("Unable to replace existing file: %v", err))
			}
		}
	}

	// * Store the file to google drive
	fileExt := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("ava-%d_%d%s", user.ID, time.Now().UnixNano(), fileExt)

	driveFile := gdstorage.StoreFileInput{Name: fileName, FileHeader: file}

	driveFileID, err := s.gds.StoreFile(&driveFile, os.Getenv("DRIVE_APP_USER_IMAGES_DIR_ID"))
	if err != nil {
		return "", err
	}

	// * Update avatar to db
	user.Avatar = driveFileID
	_, err = s.repository.Update(user)

	if err != nil {
		return "", err
	}

	return gdstorage.GetURL(driveFileID), nil

}
