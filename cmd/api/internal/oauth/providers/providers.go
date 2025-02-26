package providers

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/arinji2/vocab-thing/internal/models"
	"github.com/arinji2/vocab-thing/internal/utils/idgen"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

type ProviderInterface interface {
	FetchAuthUser(token *oauth2.Token) (*models.AuthUser, error)
	GenerateCodeURL(r *http.Request, w http.ResponseWriter) (string, error)
}

var (
	ValidProviders []string = []string{"google", "github", "discord"}
	store                   = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
)

type BaseProvider struct {
	ProviderType string
	Ctx          context.Context
	ClientId     string
	ClientSecret string
	DisplayName  string
	RedirectURL  string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
	Scopes       []string
}

func (p *BaseProvider) NewProvider(providerType string) ProviderInterface {
	if p.Ctx == nil {
		p.Ctx = context.Background()
	}
	switch providerType {
	case "google":
		return NewGoogleProvider(p.Ctx)
	case "github":
		return NewGithubProvider(p.Ctx)
	case "discord":
		return NewDiscordProvider(p.Ctx)
	default:
		return nil
	}
}

func (p *BaseProvider) FetchRawUserInfo(token *oauth2.Token) ([]byte, error) {
	req, err := http.NewRequestWithContext(p.Ctx, "GET", p.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	return p.sendRawUserInfoRequest(req, token)
}

func (p *BaseProvider) Client(token *oauth2.Token) *http.Client {
	return p.oauth2Config().Client(p.Ctx, token)
}

// sendRawUserInfoRequest sends the specified user info request and return its raw response body.
func (p *BaseProvider) sendRawUserInfoRequest(req *http.Request, token *oauth2.Token) ([]byte, error) {
	client := p.Client(token)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// http.Client.Get doesn't treat non 2xx responses as error
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf(
			"failed to fetch OAuth2 user profile via %s (%d):\n%s",
			p.UserInfoURL,
			res.StatusCode,
			string(result),
		)
	}

	return result, nil
}

// oauth2Config constructs a oauth2.Config instance based on the provider settings.
func (p *BaseProvider) oauth2Config() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  p.RedirectURL,
		ClientID:     p.ClientId,
		ClientSecret: p.ClientSecret,
		Scopes:       p.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  p.AuthURL,
			TokenURL: p.TokenURL,
		},
	}
}

func (p *BaseProvider) GenerateCodeURL(r *http.Request, w http.ResponseWriter) (string, error) {
	stateID, err := idgen.GenerateRandomID(idgen.DefaultIDSize, idgen.URLSafeAlphanumericCharset)
	if err != nil {
		return "", fmt.Errorf("error with generating state id: %w", err)
	}

	codeVerifier, err := idgen.GenerateRandomID(idgen.OauthCodeVerifierSize, idgen.NumberCharset)
	if err != nil {
		return "", fmt.Errorf("error with generating code verifier: %w", err)
	}

	hashedVerifier := sha256.Sum256([]byte(codeVerifier))
	initalCodeVerifier := base64.RawURLEncoding.EncodeToString(hashedVerifier[:])

	session, err := store.New(r, "oauth-session")
	if err != nil {
		return "", fmt.Errorf("error with creating session: %w", err)
	}
	session.Values["state"] = stateID
	session.Values["code_verifier"] = codeVerifier
	session.Save(r, w)

	codeURL := fmt.Sprintf("%s?code_challenge=%s&code_challenge_method=S256&state=%s", p.AuthURL, initalCodeVerifier, stateID)
	return codeURL, nil
}
