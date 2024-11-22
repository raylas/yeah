package main

import (
	"log"

	"github.com/alexflint/go-arg"

	"rymnd.net/yeah/internal/oui"
	"rymnd.net/yeah/internal/output"
)

type args struct {
	Mac    string `arg:"positional,required" help:"MAC address to resolve"`
	Update bool   `arg:"-u,--" help:"update the vendors database"`
	Format string `arg:"-f,--" help:"output format (table, json)" default:"table"`
}

func main() {
	var args args
	arg.MustParse(&args)

	oui, close, err := oui.New(args.Update)
	if err != nil {
		log.Fatalf("failed to initialize OUI: %v", err)
	}
	defer close()

	vendors, _, err := oui.Resolve(args.Mac)
	if err != nil {
		log.Fatalf("failed to resolve vendor: %v", err)
	}

	w, err := output.NewWriter(args.Format)
	if err != nil {
		log.Fatalf("failed to create output writer: %v", err)
	}

	if err := w.WriteHeader([]string{"OUI", "Organization"}); err != nil {
		log.Fatalf("failed to write header: %v", err)
	}

	for _, vendor := range vendors {
		if err := w.WriteResource([]output.Field{
			{T: vendor.Oui},
			{T: vendor.Org},
		}); err != nil {
			log.Fatalf("failed to write resource: %v", err)
		}
	}

	if err := w.Flush(); err != nil {
		log.Fatalf("failed to write to console: %v", err)
	}
}
