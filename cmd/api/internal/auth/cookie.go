package auth

import (
	"net/http"
	"os"
	"time"
)

func CreateUserSessionCookie(w http.ResponseWriter, sessionID string, expiresAt time.Time) {
	if expiresAt.IsZero() {
		expiresAt = time.Now().Add(time.Hour * 24 * 7)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
		Expires:  expiresAt.UTC(),
	})
}
