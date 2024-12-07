package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/alexflint/go-arg"

	"rymnd.net/yeah/internal/api"
	"rymnd.net/yeah/internal/cli"
	"rymnd.net/yeah/internal/decode"
)

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	vendors, err := decode.Vendors()
	if err != nil {
		return fmt.Errorf("failed to decode vendors: %w", err)
	}

	var args cli.Args
	arg.MustParse(&args)

	switch {
	case args.Listen:
		return api.Run(ctx, args, vendors)
	default:
		return cli.Run(ctx, args, vendors)
	}
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		panic(err)
	}
}
