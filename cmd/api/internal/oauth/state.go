package oauth

import (
	"net/http"
	"os"

	"github.com/arinji2/vocab-thing/internal/utils/idgen"
)

// GenerateState creates a random state string and stores it in a cookie
func GenerateState(r *http.Request, w http.ResponseWriter) string {
	state := idgen.GenerateRandomID(idgen.DefaultIDSize, idgen.URLSafeAlphanumericCharset)

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		SameSite: http.SameSiteStrictMode,
	})

	return state
}

// ValidateState checks if the state from the request matches the stored state
func ValidateState(r *http.Request, state string) bool {
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		return false
	}
	return cookie.Value == state
}
