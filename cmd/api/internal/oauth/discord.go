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
		Id       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := json.Unmarshal(body, &extracted); err != nil {
		return nil, err
	}

	user := &models.User{
		Username: extracted.Username,
		Email:    extracted.Email,
	}

	if o.ProviderUserID == "" {
		o.ProviderUserID = extracted.Id
	}

	return user, nil
}
