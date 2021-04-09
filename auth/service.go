package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type service struct{}

var SECRET_KEY []byte = []byte("SuperDuperUltraSecretKey")

func NewService() Service {
	return &service{}
}

func (s *service) GenerateToken(userID int) (string, error) {
	// Mapping claims (userID, iat)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
		"iat":     time.Now(),
	}

	// Generate token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sigining token
	signedToken, err := token.SignedString(SECRET_KEY)

	if err != nil {
		return signedToken, err
	}

	// Return signed token
	return signedToken, nil
}

func (s *service) ValidateToken(accessToken string) (*jwt.Token, error) {
	/*
		About jwt.Parse(token, keyFunc):
			Parse, validate, and return a token.
			keyFunc will receive the parsed token and should return the SECRET KEY (and error) for validating.
			keyFunc is an anonymous function which checks the accessToken hashing method.
			If everything is kosher, err will be nil
	*/
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// Data type assertion, it checks if the data type is the same as the given argument
		// It returns the variable and a boolean
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Invalid access token.")
		}

		// If the sigining method matches, return the SECRET KEY in byte
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return token, err
	}

	return token, nil
}
