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
	return mux
}
