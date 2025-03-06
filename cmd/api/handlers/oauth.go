package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/arinji2/vocab-thing/internal/oauth"
	"github.com/davecgh/go-spew/spew"
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
	spew.Dump(user)

	w.WriteHeader(http.StatusOK)
	writeJSON(w, http.StatusOK, p)
}
