package types

import "og-style/models"

type SignInResponse struct {
	User         *models.User `json:"user"`
	AccessToken  string       `json:"-"`
	RefreshToken string       `json:"-"`
}
