package httpmiddleware

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/arinji2/vocab-thing/internal/auth"
	"github.com/arinji2/vocab-thing/internal/database"
	"github.com/arinji2/vocab-thing/internal/errorcode"
)

func Authentication(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionID, err := auth.GetUserSession(r)
			if err != nil {
				errorcode.WriteJSONError(w, err, http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			sessionModel := database.SessionModel{DB: db}
			sessionData, err := sessionModel.Validate(ctx, sessionID)
			if err != nil {
				if err == auth.ErrSessionExpired {
					auth.DeleteUserSessionCookie(w)
				}
				fmt.Println(err.Error())
				errorcode.WriteJSONError(w, err, http.StatusUnauthorized)
				return
			}

			ctx = auth.ContextWithSession(ctx, sessionData)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
