package routes

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/arinji2/vocab-thing/handlers"
	"github.com/arinji2/vocab-thing/internal/httpmiddleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterRoutes(db *sql.DB) http.Handler {
	handler := handlers.NewHandler(db)
	userHandler := handlers.UserHandler{Handler: handler}
	phraseHandler := handlers.PhraseHandler{Handler: handler}

	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(httpmiddleware.Cors)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Timeout(60 * time.Second))

	r.Group(func(r chi.Router) {
		r.Get("/", userHandler.GetAllUsers)
		r.Post("/user/create", userHandler.CreateUser)
		r.Post("/oauth/generate-code-url", userHandler.GenerateCodeURL)
		r.Post("/oauth/callback", userHandler.CallbackHandler)
		r.Post("/user/create/guest", userHandler.CreateGuestUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(httpmiddleware.Authentication(db))
		r.Get("/user/authenticated", userHandler.AuthenticatedRoute)

		r.Route("/phrase", func(r chi.Router) {
			r.Route("/create", func(r chi.Router) {
				r.Post("/phrase", phraseHandler.CreatePhrase)
				r.Post("/tag", phraseHandler.CreateTag)
			})
			r.Get("/{id}", phraseHandler.GetPhraseByID)

			r.With(httpmiddleware.Paginate).Get("/", phraseHandler.GetAllPhrases)
			r.With(httpmiddleware.Searching).Get("/search", phraseHandler.SearchPhrases)
		})
	})

	return r
}
