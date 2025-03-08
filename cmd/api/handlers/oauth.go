package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/arinji2/vocab-thing/internal/database"
	"github.com/arinji2/vocab-thing/internal/models"
	"github.com/arinji2/vocab-thing/internal/oauth"
)

type generateCodeURLRequest struct {
	ProviderType string `json:"providerType"`
}

func (h *UserHandler) GenerateCodeURL(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var data generateCodeURLRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	provider, err := oauth.NewProvider(ctx, data.ProviderType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	codeURL, err := provider.GenerateCodeURL(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(codeURL))
}

type callbackHandlerRequest struct {
	ProviderType string `json:"providerType"`
	Code         string `json:"code"`
	State        string `json:"state"`
	Fingerprint  string `json:"fingerprint"`
	IP           string `json:"ip"`
}
type callbackHandlerResponse struct {
	SessionID string `json:"sessionID"`
}

func (h *UserHandler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var data callbackHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	provider, err := oauth.NewProvider(ctx, data.ProviderType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := provider.AuthenticateWithCode(r, data.Code, data.State)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := provider.FetchAuthUser(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userModel := database.UserModel{DB: h.DB}
	dbUser, err := userModel.ByEmail(ctx, user.Email)
	if err != nil {
		err = userModel.Create(ctx, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		user = &dbUser
	}

	providerModel := database.ProviderModel{DB: h.DB}
	selectedUserProvider := models.OauthProvider{}

	userProviders, err := providerModel.ByUserID(ctx, user.ID)
	if err != nil {
		selectedUserProvider = models.OauthProvider{
			UserID:       user.ID,
			Type:         p.Type,
			AccessToken:  p.AccessToken,
			RefreshToken: p.RefreshToken,
			ExpiresAt:    p.ExpiresAt,
		}
		err = providerModel.Create(ctx, &selectedUserProvider)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		providerExists := slices.ContainsFunc(userProviders, func(provider models.OauthProvider) bool {
			if provider.Type == p.Type {
				selectedUserProvider = provider
				return true
			}
			return false
		})

		if !providerExists {
			selectedUserProvider = models.OauthProvider{
				UserID:       user.ID,
				Type:         p.Type,
				AccessToken:  p.AccessToken,
				RefreshToken: p.RefreshToken,
				ExpiresAt:    p.ExpiresAt,
			}
			err = providerModel.Create(ctx, &selectedUserProvider)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	sessionModel := database.SessionModel{DB: h.DB}
	var userSession models.Session
	existingSessions, err := sessionModel.ByUserIDWithProvider(ctx, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, session := range existingSessions {
		if session.ProviderType == nil {
			continue
		}
		if *session.ProviderType == selectedUserProvider.Type && session.Fingerprint == data.Fingerprint && session.IP == data.IP && session.ExpiresAt.After(time.Now().Add(time.Hour*24)) {
			userSession = session
			break
		}
	}

	if userSession.ID == "" {
		userSession = models.Session{
			UserID:      user.ID,
			ProviderID:  selectedUserProvider.ID,
			Fingerprint: data.Fingerprint,
			IP:          data.IP,
			ExpiresAt:   oauth.SessionExpiry(time.Now()),
		}
		err = sessionModel.Create(ctx, &userSession)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	response := callbackHandlerResponse{
		SessionID: userSession.ID,
	}
	writeJSON(w, http.StatusOK, response)
}
