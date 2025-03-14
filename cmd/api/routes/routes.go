package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/arinji2/vocab-thing/handlers"
	"github.com/arinji2/vocab-thing/internal/auth"
	"github.com/arinji2/vocab-thing/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterRoutes(db *sql.DB) http.Handler {
	handler := handlers.NewHandler(db)
	userHandler := handlers.UserHandler{Handler: handler}
	phraseHandler := handlers.PhraseHandler{Handler: handler}

	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(corsMiddleware)

	r.Group(func(r chi.Router) {
		r.Get("/", userHandler.GetAllUsers)
		r.Post("/user/create", userHandler.CreateUser)
		r.Post("/oauth/generate-code-url", userHandler.GenerateCodeURL)
		r.Post("/oauth/callback", userHandler.CallbackHandler)
		r.Post("/user/create/guest", userHandler.CreateGuestUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(authenticatedMiddleware(db))
		r.Get("/user/authenticated", userHandler.AuthenticatedRoute)

		r.Route("/phrase", func(r chi.Router) {
			r.Route("/create", func(r chi.Router) {
				r.Post("/phrase", phraseHandler.CreatePhrase)
				r.Post("/tag", phraseHandler.CreateTag)
			})
			r.Get("/{id}", phraseHandler.GetPhraseByID)

			// Paginated list of all phrases
			// r.Get("/", phraseHandler.GetAllPhrases)
		})
	})

	return r
}

func authenticatedMiddleware(db *sql.DB) func(http.Handler) http.Handler {
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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		frontendURL := os.Getenv("FRONTEND_URL")
		fmt.Printf("Allowing CORS For URL: %s \n", frontendURL)

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
