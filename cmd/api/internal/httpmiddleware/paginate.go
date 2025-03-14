package httpmiddleware

import (
	"context"
	"net/http"
	"strconv"
)

type paginationCtxKey struct{}

type Pagination struct {
	Page     int
	PageSize int
}

func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		page, err := strconv.Atoi(query.Get("page"))
		if err != nil || page < 1 {
			page = 1
		}

		pageSize, err := strconv.Atoi(query.Get("pageSize"))
		if err != nil || pageSize < 1 || pageSize > 100 {
			pageSize = 10
		}

		ctx := context.WithValue(r.Context(), paginationCtxKey{}, Pagination{
			Page:     page,
			PageSize: pageSize,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func PaginationFromContext(ctx context.Context) (Pagination, bool) {
	pagination, ok := ctx.Value(paginationCtxKey{}).(Pagination)
	return pagination, ok
}
