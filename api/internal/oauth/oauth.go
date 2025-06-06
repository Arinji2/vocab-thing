package oauth

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/arinji2/vocab-thing/internal/errorcode"
	"github.com/arinji2/vocab-thing/internal/models"
	"github.com/gorilla/sessions"
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

var (
	sessionStore            = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	ValidProviders []string = []string{"google", "github", "discord"}
)

func NewProvider(ctx context.Context, providerType string) (ProviderInterface, error) {
	switch providerType {
	case "google":
		return NewGoogleProvider(ctx), nil
	case "github":
		return NewGithubProvider(ctx), nil
	case "discord":
		return NewDiscordProvider(ctx), nil
	default:
		return nil, errorcode.ErrUnsupportedProvider
	}
}

func (p *BaseProvider) GenerateCodeURL(r *http.Request, w http.ResponseWriter) (string, error) {
	state := GenerateState(r, w)
	session, err := sessionStore.Get(r, "oauth_session")
	if err != nil {
		log.Printf("error getting session store: %s", err.Error())
		return "", errorcode.ErrGettingSessionStore
	}
	session.Values["oauth_state"] = state

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   SessionExpiry(time.Now()).UTC().Second(),
		HttpOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		SameSite: http.SameSiteLaxMode,
	}
	err = session.Save(r, w)
	if err != nil {
		log.Printf("error saving session store: %s", err.Error())
		return "", errorcode.ErrSavingSessionStore
	}
	return p.Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce), nil
}

func (p *BaseProvider) AuthenticateWithCode(r *http.Request, code string, state string) (*models.OauthProvider, error) {
	session, err := sessionStore.Get(r, "oauth_session")
	if err != nil {
		log.Printf("error getting session store: %s", err.Error())
		return nil, errorcode.ErrGettingSessionStore
	}
	val := session.Values["oauth_state"]
	sessionState, ok := val.(string)
	if !ok {
		return nil, errorcode.ErrInvalidOauthState
	}
	state, err = url.QueryUnescape(state)
	if err != nil {
		log.Printf("error unescaping state: %s", err.Error())
		return nil, errorcode.ErrURLUnescape
	}

	code, err = url.QueryUnescape(code)
	if err != nil {
		log.Printf("error unescaping code: %s", err.Error())
		return nil, errorcode.ErrURLUnescape
	}

	if sessionState != state {
		return nil, errorcode.ErrInvalidOauthState
	}
	token, err := p.Config.Exchange(p.Ctx, code)
	if err != nil {
		log.Printf("error exchanging token: %s", err.Error())
		return nil, errorcode.ErrExchangeToken
	}
	return &models.OauthProvider{
		Type:         p.ProviderType,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
	}, nil
}

func (p *BaseProvider) RefreshAccessToken(o *models.OauthProvider) error {
	if p.ProviderType == "github" {
		return nil
	}
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
		log.Printf("error refreshing token: %s", err.Error())
		return errorcode.ErrRefreshToken
	}

	o.AccessToken = newToken.AccessToken
	o.RefreshToken = newToken.RefreshToken
	o.ExpiresAt = newToken.Expiry

	return nil
}

func SessionExpiry(t time.Time) time.Time {
	return t.Add(time.Hour * 24 * 7) // 7 days
}
