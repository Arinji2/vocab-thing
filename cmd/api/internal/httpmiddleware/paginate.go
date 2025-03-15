package httpmiddleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

type paginationCtxKey struct{}

type Pagination struct {
	Page     int
	PageSize int
	Sorting  Sorting
}

type Sorting struct {
	SortBy  string // createdAt, usageCount
	Order   string // ASC or DESC
	GroupBy string // foundIn
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

		sortBy := validateSortBy(query.Get("sortBy"))
		order := validateSortOrder(query.Get("order"))
		groupBy := validateGroupBy(query.Get("groupBy"))

		ctx := context.WithValue(r.Context(), paginationCtxKey{}, Pagination{
			Page:     page,
			PageSize: pageSize,
			Sorting: Sorting{
				SortBy:  sortBy,
				Order:   order,
				GroupBy: groupBy,
			},
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func PaginationFromContext(ctx context.Context) (Pagination, bool) {
	pagination, ok := ctx.Value(paginationCtxKey{}).(Pagination)
	return pagination, ok
}

func validateSortBy(sortBy string) string {
	switch strings.ToLower(sortBy) {
	case "createdat", "usagecount":
		return sortBy
	default:
		return "createdAt"
	}
}

func validateSortOrder(order string) string {
	switch strings.ToUpper(order) {
	case "ASC", "DESC":
		return order
	default:
		return "DESC"
	}
}

func validateGroupBy(groupBy string) string {
	if strings.ToLower(groupBy) == "foundin" {
		return groupBy
	}
	return ""
}
