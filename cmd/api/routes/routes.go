package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/arinji2/vocab-thing/handlers"
	"github.com/arinji2/vocab-thing/internal/auth"
	"github.com/arinji2/vocab-thing/internal/database"
)

func RegisterRoutes(db *sql.DB) http.Handler {
	handler := handlers.NewHandler(db)
	userHandler := handlers.UserHandler{Handler: handler}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", userHandler.GetAllUsers)
	mux.HandleFunc("POST /user/create", userHandler.CreateUser)
	mux.HandleFunc("POST /oauth/generate-code-url", userHandler.GenerateCodeURL)
	mux.HandleFunc("POST /oauth/callback", userHandler.CallbackHandler)
	mux.HandleFunc("POST /user/create/guest", userHandler.CreateGuestUser)
	mux.Handle("GET /user/authenticated", authenticatedMiddleware(
		http.HandlerFunc(userHandler.AuthenticatedRoute), db))

	return corsMiddleware(mux)
}

func authenticatedMiddleware(next http.Handler, db *sql.DB) http.Handler {
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

func corsMiddleware(next http.Handler) http.Handler {
	frontendURL := os.Getenv("FRONTEND_URL")
	fmt.Printf("Allowing CORS For URL: %s", frontendURL)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", frontendURL) // update with frontendURL
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Expose-Headers", "Link")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
