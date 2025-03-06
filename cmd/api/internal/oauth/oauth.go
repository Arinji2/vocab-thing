package oauth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/arinji2/vocab-thing/internal/models"
	"golang.org/x/oauth2"
)

type ProviderInterface interface {
	FetchAuthUser(o *models.OauthProvider) (*models.User, error)
	GenerateCodeURL(r *http.Request, w http.ResponseWriter) (string, error)
	AuthenticateWithCode(r *http.Request, code, state string) (*models.OauthProvider, error)
	RefreshAccessToken(o *models.OauthProvider) error
}

type BaseProvider struct {
	ProviderType string
	Ctx          context.Context
	Config       *oauth2.Config
	UserInfoURL  string
}

var ValidProviders []string = []string{"google", "github", "discord"}

func NewProvider(ctx context.Context, providerType string) (ProviderInterface, error) {
	switch providerType {
	case "google":
		return NewGoogleProvider(ctx), nil
	case "github":
		return NewGithubProvider(ctx), nil
	case "discord":
		return NewDiscordProvider(ctx), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerType)
	}
}

func (p *BaseProvider) GenerateCodeURL(r *http.Request, w http.ResponseWriter) (string, error) {
	state := GenerateState(r, w)
	return p.Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce), nil
}

func (p *BaseProvider) AuthenticateWithCode(r *http.Request, code, state string) (*models.OauthProvider, error) {
	if !ValidateState(r, state) {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := p.Config.Exchange(p.Ctx, code)
	if err != nil {
		return nil, fmt.Errorf("error exchanging token: %w", err)
	}
	return &models.OauthProvider{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
	}, nil
}

func (p *BaseProvider) RefreshAccessToken(o *models.OauthProvider) error {
	existingToken := &oauth2.Token{
		AccessToken:  o.AccessToken,
		RefreshToken: o.RefreshToken,
		Expiry:       o.ExpiresAt,
	}

	// Check if the token is expired or about to expire (within 5 minutes)
	if existingToken.Valid() && time.Until(existingToken.Expiry) > 5*time.Minute {
		return nil
	}

	tokenSource := p.Config.TokenSource(p.Ctx, existingToken)
	newToken, err := tokenSource.Token()
	if err != nil {
		return fmt.Errorf("failed to refresh access token: %w", err)
	}

	o.AccessToken = newToken.AccessToken
	o.RefreshToken = newToken.RefreshToken
	o.ExpiresAt = newToken.Expiry

	return nil
}

func SessionExpiry(t time.Time) time.Time {
	return t.Add(time.Hour * 24 * 7) // 7 days
}
