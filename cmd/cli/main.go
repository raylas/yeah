package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"

	"rymnd.net/yeah/internal/decode"
	"rymnd.net/yeah/internal/output"
	"rymnd.net/yeah/internal/vendors"
)

type args struct {
	Macs   []string `arg:"positional,required"`
	Wide   bool     `arg:"-w,--" help:"include additional fields"`
	Output string   `arg:"-o,--" help:"output format (table,json)" default:"table"`
}

func main() {
	var args args
	arg.MustParse(&args)

	v, err := decode.Vendors()
	if err != nil {
		fmt.Printf("failed to load vendors: %v\n", err)
		os.Exit(1)
	}

	var results []*vendors.VendorEntry
	for _, mac := range args.Macs {
		results = append(results, v.Search(mac)...)
	}

	w, err := output.NewWriter(args.Output)
	if err != nil {
		fmt.Printf("failed to create output writer: %v\n", err)
		os.Exit(1)
	}

	headers := []string{"OUI", "Organization"}
	if args.Wide {
		headers = append(headers, "Address")
	}

	if err := w.WriteHeader(headers); err != nil {
		fmt.Printf("failed to write header: %v\n", err)
		os.Exit(1)
	}

	for _, vendor := range results {
		fields := []output.Field{
			{T: vendor.Oui},
			{T: vendor.Name},
		}
		if args.Wide {
			fields = append(fields, output.Field{T: vendor.Address})
		}

		if err := w.WriteResource(fields); err != nil {
			fmt.Printf("failed to write resource: %v\n", err)
			os.Exit(1)
		}
	}

	if err := w.Flush(); err != nil {
		fmt.Printf("failed to write to console: %v\n", err)
		os.Exit(1)
	}
}
