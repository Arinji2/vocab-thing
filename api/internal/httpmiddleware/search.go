package httpmiddleware

import (
	"context"
	"net/http"

	"github.com/arinji2/vocab-thing/internal/errorcode"
)

type searchCtxKey struct{}

type Search struct {
	Term string
}

func Searching(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		searchTerm := query.Get("searchTerm")

		if searchTerm == "" {
			errorcode.WriteJSONError(w, errorcode.ErrNoSearchingData, http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), searchCtxKey{}, Search{
			Term: searchTerm,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SearchingFromContext(ctx context.Context) (Search, bool) {
	searching, ok := ctx.Value(searchCtxKey{}).(Search)
	return searching, ok
}
