package auth

import (
	"context"
	"errors"

	"github.com/arinji2/vocab-thing/internal/models"
)

type sessionCtxKey struct{}

var ErrSessionExpired = errors.New("session has expired")

func ContextWithSession(ctx context.Context, session models.Session) context.Context {
	return context.WithValue(ctx, sessionCtxKey{}, session)
}

func SessionFromContext(ctx context.Context) (models.Session, bool) {
	session, ok := ctx.Value(sessionCtxKey{}).(models.Session)
	return session, ok
}
