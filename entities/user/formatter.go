package user

import "github.com/muktiwbw/gdstorage"

type UserFormat struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
	Token      string `json:"token"`
	Avatar     string `json:"avatar"`
}

func FormatUser(user User, token string) UserFormat {
	var avatar string

	if user.Avatar != "" {
		avatar = gdstorage.GetURL(avatar)
	}

	return UserFormat{
		ID:         user.ID,
		Name:       user.Name,
		Occupation: user.Occupation,
		Email:      user.Email,
		Avatar:     avatar,
		Token:      token,
	}
}
