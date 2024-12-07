package encode

import (
	"context"
	"encoding/gob"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"rymnd.net/yeah/internal/vendors"
)

const (
	vendorsDest = "internal/data/vendors.gob"
)

func Vendors(ctx context.Context) error {
	v := vendors.New()
	if err := collect(ctx, v); err != nil {
		return err
	}
	if err := serialize(v, vendorsDest); err != nil {
		return err
	}
	log.Ctx(ctx).Info().Str("dest", vendorsDest).Msg("written")
	return nil
}

func serialize(v *vendors.Vendors, d string) error {
	if err := os.MkdirAll(filepath.Dir(d), 0755); err != nil {
		return err
	}

	file, err := os.Create(d)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(v)
}
