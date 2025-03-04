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

type Github struct {
	BaseProvider
}

func NewGithubProvider(ctx context.Context) *Github {
	return &Github{
		BaseProvider{
			ProviderType: "github",
			Ctx:          ctx,
			Config: &oauth2.Config{
				ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
				ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
				RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
				Scopes:       []string{"read:user", "user:email"},
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://github.com/login/oauth/authorize",
					TokenURL: "https://github.com/login/oauth/access_token",
				},
			},
			UserInfoURL: "https://api.github.com/user",
		},
	}
}

func (p *Github) FetchAuthUser(o *models.OauthProvider) (*models.User, error) {
	err := p.RefreshAccessToken(o)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(p.Ctx, "GET", p.UserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.AccessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

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
		Login string `json:"login"`
		Name  string `json:"name"`
		Email string `json:"email"`
		ID    int64  `json:"id"`
	}

	if err := json.Unmarshal(body, &extracted); err != nil {
		return nil, err
	}

	// If primary email is not returned, fetch it
	if extracted.Email == "" {
		email, err := p.fetchPrimaryEmail(o)
		if err != nil {
			return nil, err
		}
		extracted.Email = email
	}

	user := &models.User{
		Username: extracted.Login,
		Email:    extracted.Email,
	}

	if o.ProviderUserID == "" {
		o.ProviderUserID = fmt.Sprintf("%d", extracted.ID)
	}

	return user, nil
}

// fetchPrimaryEmail retrieves the user's primary verified email
func (p *Github) fetchPrimaryEmail(o *models.OauthProvider) (string, error) {
	req, err := http.NewRequestWithContext(p.Ctx, "GET", p.UserInfoURL+"/emails", nil)
	if err != nil {
		return "", fmt.Errorf("error creating email request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.AccessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making email request: %w", err)
	}
	defer resp.Body.Close()

	// Ignore common authorization errors
	if resp.StatusCode == http.StatusUnauthorized ||
		resp.StatusCode == http.StatusForbidden ||
		resp.StatusCode == http.StatusNotFound {
		return "", nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read email response body: %w", err)
	}

	var emails []struct {
		Email    string `json:"email"`
		Verified bool   `json:"verified"`
		Primary  bool   `json:"primary"`
	}

	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}

	// Find first verified primary email
	for _, email := range emails {
		if email.Verified && email.Primary {
			return email.Email, nil
		}
	}

	return "", nil
}
