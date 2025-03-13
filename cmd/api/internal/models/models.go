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
	ID           string    `json:"id" sql:"id"`
	UserID       string    `json:"user_id" sql:"userId"`
	ProviderID   string    `json:"provider_id" sql:"providerId"`
	ProviderType *string   `json:"provider_type,omitempty" sql:"type"`
	Fingerprint  string    `json:"fingerprint" sql:"fingerprint"`
	IP           string    `json:"ip" sql:"ip"`
	ExpiresAt    time.Time `json:"expires_at" sql:"expiresAt"`
	CreatedAt    time.Time `json:"created_at" sql:"createdAt"`
}

type Phrase struct {
	ID               string    `json:"id" sql:"id"`
	UserID           string    `json:"user_id" sql:"userId"`
	Phrase           string    `json:"phrase" sql:"phrase"`
	PhraseDefinition string    `json:"phrase_definition" sql:"phraseDefinition"`
	Pinned           bool      `json:"pinned" sql:"pinned"`
	FoundIn          string    `json:"found_in" sql:"foundIn"`
	Public           bool      `json:"public" sql:"public"`
	UsageCount       int       `json:"usage_count" sql:"usageCount"`
	CreatedAt        time.Time `json:"created_at" sql:"createdAt"`
}

type PhraseTag struct {
	ID        string    `json:"id" sql:"id"`
	PhraseID  string    `json:"phrase_id" sql:"phraseId"`
	TagName   string    `json:"tag_name" sql:"tagName"`
	TagColor  string    `json:"tag_color" sql:"tagColor"`
	CreatedAt time.Time `json:"created_at" sql:"createdAt"`
}

type TaggedPhrase struct {
	Phrase Phrase
	Tag    []PhraseTag
}
