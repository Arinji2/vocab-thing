package providers

import (
	"context"
	"encoding/json"
	"os"

	"github.com/arinji2/vocab-thing/internal/tools/types"
	"golang.org/x/oauth2"
)

type Google struct {
	BaseProvider
}

func NewGoogleProvider() *Google {
	return &Google{
		BaseProvider{
			ProviderType: "google",
			Ctx:          context.Background(),
			AuthURL:      "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL:     "https://oauth2.googleapis.com/token",
			UserInfoURL:  "https://www.googleapis.com/oauth2/v3/userinfo",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email",
			},
			ClientId:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		},
	}
}

// FetchAuthUser returns an AuthUser instance based on the Google's user api.
func (p *Google) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id            string `json:"sub"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Type:         p.ProviderType,
		Id:           extracted.Id,
		Username:     extracted.Name,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	if extracted.EmailVerified {
		user.Email = extracted.Email
	}

	return user, nil
}
