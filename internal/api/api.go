package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"rymnd.net/yeah/internal/cli"
	"rymnd.net/yeah/internal/vendors"
)

func Run(ctx context.Context, args cli.Args, v *vendors.Vendors) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Set up logging
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp}
	level, _ := zerolog.ParseLevel(args.LogLevel)
	zerolog.SetGlobalLevel(level)
	ctx = zerolog.New(output).With().
		Timestamp().
		Logger().
		WithContext(ctx)

	// Set up routes
	m := http.NewServeMux()
	m.Handle("/{query}", handleSearch(v))

	// Start the HTTP server
	s := new(ctx, m)
	if err := serve(ctx, args.Bind, s); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func new(ctx context.Context, mux *http.ServeMux) http.Handler {
	var handler http.Handler = mux
	handler = logRequests(ctx, handler)
	return handler
}

func serve(ctx context.Context, listenAddress string, server http.Handler) error {
	httpServer := &http.Server{
		Addr:    listenAddress,
		Handler: server,
	}

	go func() {
		log.Ctx(ctx).Info().Str("address", httpServer.Addr).Msg("server listening")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Ctx(ctx).Error().Err(err).Msg("server error")
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Ctx(ctx).Info().Msg("server shutting down")

		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("server shutdown error")
		}
	}()
	wg.Wait()

	return nil
}

func logRequests(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Ctx(ctx).Debug().Str("method", r.Method).Str("path", r.URL.Path).Msg("request")
		next.ServeHTTP(w, r)
	})
}
