package encode

import (
	"encoding/gob"
	"os"
	"path/filepath"

	"rymnd.net/yeah/internal/vendors"
)

const (
	vendorsDest = "internal/data/vendors.gob"
)

func Vendors() error {
	v := vendors.New()
	if err := collect(v); err != nil {
		return err
	}
	return serialize(v, vendorsDest)
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
