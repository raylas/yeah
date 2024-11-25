package encode

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"rymnd.net/yeah/internal/vendors"
)

var (
	sources = []string{
		"https://standards-oui.ieee.org/oui/oui.csv",
		"https://standards-oui.ieee.org/cid/cid.csv",
		"https://standards-oui.ieee.org/iab/iab.csv",
		"https://standards-oui.ieee.org/oui28/mam.csv",
		"https://standards-oui.ieee.org/oui36/oui36.csv",
	}
)

func collect(ctx context.Context, v *vendors.Vendors) error {
	for _, source := range sources {
		if err := download(v, source); err != nil {
			return err
		}
		log.Ctx(ctx).Info().Str("source", source).Msg("encoded")
	}
	return nil
}

func download(v *vendors.Vendors, source string) error {
	resp, err := http.Get(source)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download %s: %s", source, resp.Status)
	}

	// Reference the source
	if err := v.Reference(source, resp.Header.Get("Last-Modified")); err != nil {
		return fmt.Errorf("failed to reference source: %w", err)
	}

	// Parse the CSV
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		// Skip "incomplete" records
		if len(record) < 3 {
			continue
		}

		oui := record[1]
		org := record[2]
		address := record[3]

		v.Insert(oui, &vendors.VendorEntry{
			Oui:     oui,
			Name:    org,
			Address: address,
		})
	}

	return nil
}
