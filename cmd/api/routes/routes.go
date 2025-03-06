package routes

import (
	"database/sql"
	"net/http"

	"github.com/arinji2/vocab-thing/handlers"
)

func RegisterRoutes(db *sql.DB) http.Handler {
	handler := handlers.NewHandler(db)
	userHandler := handlers.UserHandler{Handler: handler}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", userHandler.GetAllUsers)
	mux.HandleFunc("POST /user/create", userHandler.CreateUser)
	mux.HandleFunc("POST /oauth/generate-code-url", userHandler.GenerateCodeURL)
	return mux
	mux.HandleFunc("POST /oauth/callback", userHandler.CallbackHandler)
}
