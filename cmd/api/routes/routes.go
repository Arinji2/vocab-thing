package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/arinji2/vocab-thing/handlers"
)

func RegisterRoutes(db *sql.DB) http.Handler {
	handler := handlers.NewHandler(db)
	userHandler := handlers.UserHandler{Handler: handler}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", userHandler.GetAllUsers)
	mux.HandleFunc("POST /user/create", userHandler.CreateUser)
	mux.HandleFunc("POST /oauth/generate-code-url", userHandler.GenerateCodeURL)
	mux.HandleFunc("POST /oauth/callback", userHandler.CallbackHandler)
	return corsMiddleware(mux)
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
