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

func DeleteUserSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Now().UTC(),
	})
}

func GetUserSession(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
