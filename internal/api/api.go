package api

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

func New(ctx context.Context, mux *http.ServeMux) http.Handler {
	var handler http.Handler = mux
	handler = logRequests(ctx, handler)
	return handler
}

func Serve(ctx context.Context, listenAddress string, server http.Handler) error {
	httpServer := &http.Server{
		Addr:    listenAddress,
		Handler: server,
	}

	go func() {
		log.Ctx(ctx).Info().Str("address", httpServer.Addr).Msg("server listening")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Ctx(ctx).Error().Err(err).Msg("api error")
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("api shutdown error")
		}
	}()
	wg.Wait()

	return nil
}

func handleError(err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
	})
}
