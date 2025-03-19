package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/arinji2/vocab-thing/internal/errorcode"
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
		log.Printf("error refreshing access token: %s", err.Error())
		return nil, errorcode.ErrRefreshToken
	}

	req, err := http.NewRequestWithContext(p.Ctx, "GET", p.UserInfoURL, nil)
	if err != nil {
		log.Printf("error creating request: %s", err.Error())
		return nil, errorcode.ErrFetchingOauthUser
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.AccessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error making request: %s", err.Error())
		return nil, errorcode.ErrFetchingOauthUser
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("unexpected status code: %d", resp.StatusCode)
		return nil, errorcode.ErrFetchingOauthUser
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err.Error())
		return nil, errorcode.ErrFetchingOauthUser
	}

	var extracted struct {
		Login string `json:"login"`
		Name  string `json:"name"`
		Email string `json:"email"`
		ID    int64  `json:"id"`
	}

	if err := json.Unmarshal(body, &extracted); err != nil {
		log.Printf("error unmarshalling response body: %s", err.Error())
		return nil, errorcode.ErrFetchingOauthUser
	}

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

	if o.UserID == "" {
		o.UserID = fmt.Sprintf("%d", extracted.ID)
	}

	return user, nil
}

func (p *Github) fetchPrimaryEmail(o *models.OauthProvider) (string, error) {
	req, err := http.NewRequestWithContext(p.Ctx, "GET", p.UserInfoURL+"/emails", nil)
	if err != nil {
		log.Printf("error creating email request: %s", err.Error())
		return "", errorcode.ErrFetchingOauthUser
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.AccessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error making email request: %s", err.Error())
		return "", errorcode.ErrFetchingOauthUser
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized ||
		resp.StatusCode == http.StatusForbidden ||
		resp.StatusCode == http.StatusNotFound {
		return "", nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading email response body: %s", err.Error())
		return "", errorcode.ErrFetchingOauthUser
	}

	var emails []struct {
		Email    string `json:"email"`
		Verified bool   `json:"verified"`
		Primary  bool   `json:"primary"`
	}

	if err := json.Unmarshal(body, &emails); err != nil {
		log.Printf("error unmarshalling email response body: %s", err.Error())
		return "", errorcode.ErrFetchingOauthUser
	}

	for _, email := range emails {
		if email.Verified && email.Primary {
			return email.Email, nil
		}
	}

	return "", nil
}
