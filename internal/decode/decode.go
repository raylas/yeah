package decode

import (
	"bytes"
	_ "embed"
	"encoding/gob"

	"rymnd.net/yeah/internal/data"
	"rymnd.net/yeah/internal/vendors"
)

func Vendors() (*vendors.Vendors, error) {
	var v vendors.Vendors
	if err := gob.NewDecoder(bytes.NewReader(data.Vendors)).Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}
