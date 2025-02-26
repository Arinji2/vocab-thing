package models

import "github.com/arinji2/vocab-thing/internal/utils/datetime"

type AuthUser struct {
	Type         string            `json:"type"`
	Expiry       datetime.DateTime `json:"expiry"`
	Id           string            `json:"id"`
	Username     string            `json:"username"`
	Email        string            `json:"email"`
	AccessToken  string            `json:"accessToken"`
	RefreshToken string            `json:"refreshToken"`
	IsGuest      bool              `json:"isGuest"`
}
