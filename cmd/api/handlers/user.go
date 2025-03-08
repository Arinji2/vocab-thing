package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/arinji2/vocab-thing/internal/auth"
	"github.com/arinji2/vocab-thing/internal/database"
	"github.com/arinji2/vocab-thing/internal/models"
	"github.com/arinji2/vocab-thing/internal/oauth"
	"github.com/davecgh/go-spew/spew"
)

type UserHandler struct {
	*Handler
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	userModel := database.UserModel{DB: h.DB}
	users, err := userModel.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, user := range users {
		writeJSON(w, http.StatusOK, user)
	}
}

type createUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var data createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userModel := database.UserModel{DB: h.DB}
	userData := models.User{
		Email:    data.Email,
		Username: data.Username,
	}
	err := userModel.Create(ctx, &userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	writeJSON(w, http.StatusOK, userData)
}

func (h *UserHandler) CreateGuestUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	provider := oauth.NewGuestProvider(h.DB)

	user, err := provider.FetchGuestUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	spew.Dump(user)
	userModel := database.UserModel{DB: h.DB}
	err = userModel.Create(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	providerModel := database.ProviderModel{DB: h.DB}

	userProvider := models.OauthProvider{
		UserID:       user.ID,
		Type:         "guest",
		AccessToken:  "",
		RefreshToken: "",
		ExpiresAt:    time.Time{},
	}
	err = providerModel.Create(ctx, &userProvider)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionModel := database.SessionModel{DB: h.DB}

	userSession := models.Session{
		UserID:      user.ID,
		ProviderID:  userProvider.ID,
		Fingerprint: "",
		IP:          "",
		ExpiresAt:   time.Now().Add(365 * 24 * time.Hour), // 1 year
	}
	err = sessionModel.Create(ctx, &userSession)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	auth.CreateUserSessionCookie(w, userSession.ID, userSession.ExpiresAt)
	w.WriteHeader(http.StatusOK)
}
