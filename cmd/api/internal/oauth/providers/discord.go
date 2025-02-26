package providers

import (
	"context"
	"encoding/json"
	"os"

	"github.com/arinji2/vocab-thing/internal/tools/types"
	"golang.org/x/oauth2"
)

type Discord struct {
	BaseProvider
}

func NewDiscordProvider() *Discord {
	return &Discord{
		BaseProvider{
			ProviderType: "discord",
			Ctx:          context.Background(),
			AuthURL:      "https://discord.com/oauth2/authorize",
			TokenURL:     "https://discord.com/api/oauth2/token",
			UserInfoURL:  "https://discord.com/api/users/@me",
			Scopes: []string{
				"identify",
				"email",
			},
			ClientId:     os.Getenv("DISCORD_CLIENT_ID"),
			ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("DISCORD_REDIRECT_URL"),
		},
	}
}

// FetchAuthUser returns an AuthUser instance based on the Discord's user api.
func (p *Discord) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Avatar        string `json:"avatar"`
		Email         string `json:"email"`
		Verified      bool   `json:"verified"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Username:     extracted.Username,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	if extracted.Verified {
		user.Email = extracted.Email
	}

	return user, nil
}
