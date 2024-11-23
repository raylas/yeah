package cli

import (
	"context"
	"fmt"

	"rymnd.net/yeah/internal/output"
	"rymnd.net/yeah/internal/vendors"
)

type Args struct {
	Macs     []string `arg:"positional"`
	Wide     bool     `arg:"-w,--" help:"include additional fields"`
	Output   string   `arg:"-o,--" help:"output format (table,json)" default:"table"`
	Listen   bool     `arg:"-l,--" help:"run server"`
	Bind     string   `arg:"-b,--" help:"server bind address" default:":8080"`
	LogLevel string   `arg:"-v,--" help:"log level (info,debug)" default:"info"`
}

func Run(ctx context.Context, args Args, v *vendors.Vendors) error {
	// Search for vendors
	var results []*vendors.VendorEntry
	for _, mac := range args.Macs {
		results = append(results, v.Search(mac)...)
	}

	// Create output writer
	w, err := output.NewWriter(args.Output)
	if err != nil {
		return fmt.Errorf("failed to create output writer: %w", err)
	}

	// Create output headers
	headers := []string{"OUI", "Organization"}
	if args.Wide {
		headers = append(headers, "Address")
	}

	// Write output headers
	if err := w.WriteHeader(headers); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write output resources
	for _, vendor := range results {
		fields := []output.Field{
			{T: vendor.Oui},
			{T: vendor.Name},
		}
		if args.Wide {
			fields = append(fields, output.Field{T: vendor.Address})
		}

		if err := w.WriteResource(fields); err != nil {
			return fmt.Errorf("failed to write resource: %w", err)
		}
	}

	// Flush output
	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to write to console: %w", err)
	}

	return nil
}
