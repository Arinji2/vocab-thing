package httpmiddleware

import (
	"database/sql"
	"net/http"

	"github.com/arinji2/vocab-thing/internal/auth"
	"github.com/arinji2/vocab-thing/internal/database"
)

func Authentication(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionID, err := auth.GetUserSession(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			sessionModel := database.SessionModel{DB: db}
			sessionData, err := sessionModel.Validate(ctx, sessionID)
			if err != nil {
				if err == auth.ErrSessionExpired {
					auth.DeleteUserSessionCookie(w)
				}
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx = auth.ContextWithSession(ctx, sessionData)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
