package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"rymnd.net/yeah/internal/api"
	"rymnd.net/yeah/internal/decode"
)

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Set up logging
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp}
	ctx = zerolog.New(output).With().
		Timestamp().
		Logger().
		WithContext(ctx)

	// Load vendors
	v, err := decode.Vendors()
	if err != nil {
		return fmt.Errorf("failed to load vendors: %w", err)
	}

	// Set up routes
	m := http.NewServeMux()
	m.Handle("/search", api.HandleSearch(v))

	// Start the HTTP server
	s := api.New(ctx, m)
	if err := api.Serve(ctx, ":8080", s); err != nil {
		return fmt.Errorf("unable to start server: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	log.Ctx(ctx).Info().Msg("exiting")
}
