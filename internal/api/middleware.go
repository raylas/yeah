package api

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
)

func logRequests(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Ctx(ctx).Debug().Str("method", r.Method).Str("path", r.URL.Path).Msg("request")
		next.ServeHTTP(w, r)
	})
}
