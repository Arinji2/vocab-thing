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

type Discord struct {
	BaseProvider
}

func NewDiscordProvider(ctx context.Context) *Discord {
	return &Discord{
		BaseProvider{
			ProviderType: "discord",
			Ctx:          ctx,
			Config: &oauth2.Config{
				ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
				ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
				RedirectURL:  os.Getenv("DISCORD_REDIRECT_URL"),
				Scopes:       []string{"identify", "email"},
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://discord.com/oauth2/authorize",
					TokenURL: "https://discord.com/api/oauth2/token",
				},
			},
			UserInfoURL: "https://discord.com/api/users/@me",
		},
	}
}

func (p *Discord) FetchAuthUser(o *models.OauthProvider) (*models.User, error) {
	if err := p.RefreshAccessToken(o); err != nil {
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
		Id       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := json.Unmarshal(body, &extracted); err != nil {
		log.Printf("error unmarshalling response body: %s", err.Error())
		return nil, errorcode.ErrFetchingOauthUser
	}

	user := &models.User{
		Username: extracted.Username,
		Email:    extracted.Email,
	}

	if o.UserID == "" {
		o.UserID = extracted.Id
	}

	return user, nil
}
