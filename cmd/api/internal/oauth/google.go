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
		log.Printf("error refreshing access token: %s", err.Error())
		return nil, errorcode.ErrRefreshToken
	}

	req, err := http.NewRequestWithContext(p.Ctx, "GET", p.UserInfoURL, nil)
	if err != nil {
		log.Printf("error creating request: %s", err.Error())
		return nil, errorcode.ErrFetchingOauthUser
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.AccessToken))

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
		Sub           string `json:"sub"`
		Name          string `json:"name"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}

	if err := json.Unmarshal(body, &extracted); err != nil {
		log.Printf("error unmarshalling response body: %s", err.Error())
		return nil, errorcode.ErrFetchingOauthUser
	}

	user := &models.User{
		Username: extracted.Name,
		Email:    extracted.Email,
	}

	if o.UserID == "" {
		o.UserID = extracted.Sub
	}

	return user, nil
}
