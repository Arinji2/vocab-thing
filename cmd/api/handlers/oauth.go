package handlers

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/arinji2/vocab-thing/internal/oauth/providers"
)

type generateCodeURL struct {
	ProviderType string `json:"providerType"`
}

func (h *UserHandler) GenerateCodeURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var data generateCodeURL
	err := parseRequestBody(r, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if slices.Contains(providers.ValidProviders, data.ProviderType) {
		baseProvider := providers.BaseProvider{}
		baseProvider.Ctx = ctx
		provider := baseProvider.NewProvider(data.ProviderType)
		if provider == nil {
			http.Error(w, "Invalid provider type", http.StatusBadRequest)
			return
		}
		codeURL, err := provider.GenerateCodeURL(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(codeURL))
		w.WriteHeader(http.StatusOK)
		return
	}
}
