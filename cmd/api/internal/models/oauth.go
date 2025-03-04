package models

import "time"

type OauthProvider struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"`
	UserID         string    `json:"user_id"`
	ProviderUserID string    `json:"provider_user_id"`
	AccessToken    string    `json:"access_token"`
	ExpiresIn      time.Time `json:"expires_in"`
	RefreshToken   string    `json:"refresh_token"`
	CreatedAt      string    `json:"created_at"`
}
