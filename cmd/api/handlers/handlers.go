package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/arinji2/vocab-thing/internal/auth"
	"github.com/arinji2/vocab-thing/internal/database"
)

type Handler struct {
	DB *sql.DB
}

// NewHandler creates a new base Handler.
func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

type SessionCtxKey struct{}

func authenticatedRoute(r *http.Request, db *sql.DB) (context.Context, error) {
	ctx := r.Context()
	sessionModel := database.SessionModel{DB: db}
	sessionID, err := auth.GetUserSession(r)
	if err != nil {
		return nil, err
	}

	sessionData, err := sessionModel.Validate(ctx, sessionID)
	if err != nil {
		if err == auth.ErrSessionExpired {
			return nil, auth.ErrSessionExpired
		}
		return nil, err
	}

	return auth.ContextWithSession(ctx, sessionData), nil
}
