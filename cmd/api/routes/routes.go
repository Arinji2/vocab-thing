package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/arinji2/vocab-thing/handlers"
	"github.com/arinji2/vocab-thing/internal/auth"
	"github.com/arinji2/vocab-thing/internal/database"
)

func RegisterRoutes(db *sql.DB) http.Handler {
	handler := handlers.NewHandler(db)
	userHandler := handlers.UserHandler{Handler: handler}
	phraseHandler := handlers.PhraseHandler{Handler: handler}
	mux := http.NewServeMux()

	logRoute(mux, "GET", "/", userHandler.GetAllUsers, nil)
	logRoute(mux, "POST", "/user/create", userHandler.CreateUser, nil)
	logRoute(mux, "POST", "/oauth/generate-code-url", userHandler.GenerateCodeURL, nil)
	logRoute(mux, "POST", "/oauth/callback", userHandler.CallbackHandler, nil)
	logRoute(mux, "POST", "/user/create/guest", userHandler.CreateGuestUser, nil)
	logRoute(mux, "GET", "/user/authenticated", userHandler.AuthenticatedRoute, db)
	logRoute(mux, "POST", "/phrase/create/phrase", phraseHandler.CreatePhrase, db)
	logRoute(mux, "POST", "/phrase/create/tag", phraseHandler.CreateTag, db)

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
		w.Header().Set("Access-Control-Allow-Origin", frontendURL)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Expose-Headers", "Link")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func logRoute(mux *http.ServeMux, method, route string, handlerFunc http.HandlerFunc, db *sql.DB) {
	if db != nil {
		mux.Handle(method+" "+route, authenticatedMiddleware(handlerFunc, db))
		log.Printf("Registered %s route for %s (authenticated)\n", route, method)
	} else {
		mux.HandleFunc(method+" "+route, handlerFunc)
		log.Printf("Registered %s route for %s\n", route, method)
	}
}
