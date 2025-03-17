package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/arinji2/vocab-thing/internal/auth"
	"github.com/arinji2/vocab-thing/internal/database"
)

type SyncHandler struct {
	*Handler
}

func (s *SyncHandler) GetSync(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userSession, ok := auth.SessionFromContext(ctx)
	if !ok {
		http.Error(w, "no session found", http.StatusInternalServerError)
		return
	}

	syncModel := database.SyncModel{DB: s.DB}
	responseData, err := syncModel.ByUserID(ctx, userSession.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, responseData)
}

func (s *SyncHandler) ManualSync(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userSession, ok := auth.SessionFromContext(ctx)
	if !ok {
		http.Error(w, "no session found", http.StatusInternalServerError)
		return
	}

	syncModel := database.SyncModel{DB: s.DB}
	responseData, err := syncModel.ManualSync(ctx, userSession.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, responseData)
}
