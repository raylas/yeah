package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"rymnd.net/yeah/internal/vendors"
)

func handleSearch(v *vendors.Vendors) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.PathValue("query")
		macs := strings.Split(query, ",")

		var results []*vendors.VendorEntry
		for _, mac := range macs {
			results = append(results, v.Search(mac)...)
		}

		json.NewEncoder(w).Encode(results)
	})
}
