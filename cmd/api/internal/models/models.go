package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

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

type Session struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	Fingerprint string    `json:"fingerprint"`
	IP          string    `json:"ip"`
	ExpiresAt   time.Time `json:"expires_at"`
}
