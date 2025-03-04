package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/arinji2/vocab-thing/internal/models"
	"golang.org/x/oauth2"
)

type Google struct {
	BaseProvider
}

func NewGoogleProvider(ctx context.Context) *Google {
	return &Google{
		BaseProvider{
			ProviderType: "google",
			Ctx:          ctx,
			Config: &oauth2.Config{
				ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
				ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
				RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
				Scopes:       []string{"openid", "email", "profile"},
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
					TokenURL: "https://oauth2.googleapis.com/token",
				},
			},
			UserInfoURL: "https://openidconnect.googleapis.com/v1/userinfo",
		},
	}
}

func (p *Google) FetchAuthUser(o *models.OauthProvider) (*models.User, error) {
	err := p.RefreshAccessToken(o)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(p.Ctx, "GET", p.UserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.AccessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var extracted struct {
		Sub           string `json:"sub"`
		Name          string `json:"name"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}

	if err := json.Unmarshal(body, &extracted); err != nil {
		return nil, err
	}

	user := &models.User{
		Username: extracted.Name,
		Email:    extracted.Email,
	}

	if o.ProviderUserID == "" {
		o.ProviderUserID = extracted.Sub
	}

	return user, nil
}
