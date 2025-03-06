package models

import "time"

type User struct {
	ID        string    `json:"id" sql:"id"`
	Username  string    `json:"username" sql:"username"`
	Email     string    `json:"email" sql:"email"`
	CreatedAt time.Time `json:"created_at" sql:"createdAt"`
}

type OauthProvider struct {
	ID           string    `json:"id" sql:"id"`
	UserID       string    `json:"user_id" sql:"userId"`
	Type         string    `json:"type" sql:"type"`
	RefreshToken string    `json:"refresh_token" sql:"refreshToken"`
	AccessToken  string    `json:"access_token" sql:"accessToken"`
	ExpiresAt    time.Time `json:"expires_at" sql:"expiresAt"`
	CreatedAt    time.Time `json:"created_at" sql:"createdAt"`
}

type Session struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ProviderID  string    `json:"provider_id"`
	CreatedAt   time.Time `json:"created_at"`
	Fingerprint string    `json:"fingerprint"`
	IP          string    `json:"ip"`
	ExpiresAt   time.Time `json:"expires_at"`
}
