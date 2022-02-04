package middleware

import (
	"context"
	"net/http"
	"site/internal/logger"
	"strconv"
)

func Paginate(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		var (
			page int
			err error
		)
		if pageStr != "" {
			if page, err = strconv.Atoi(pageStr); err != nil {
				logger.Logger.Error(err.Error())
			}
		}

		if page == 0 {
			page = 1
		}

		filter := r.URL.Query().Get("filter")

		ctx := r.Context()
		ctx = context.WithValue(ctx, CtxKeyPage, int32(page))
		ctx = context.WithValue(ctx, CtxKeyFilter, filter)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}