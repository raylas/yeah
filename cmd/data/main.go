package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"rymnd.net/yeah/internal/encode"
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

	if err := encode.Vendors(ctx); err != nil {
		return fmt.Errorf("failed to encode vendors: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Ctx(ctx).Fatal().Err(err)
	}
	log.Ctx(ctx).Info().Msg("exiting")
}
