package api

import (
	"encoding/json"
	"net/http"

	"rymnd.net/yeah/internal/vendors"
)

type searchRequest struct {
	Macs []string `json:"macs"`
}

func HandleSearch(v *vendors.Vendors) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
			return
		}

		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			json.NewEncoder(w).Encode(map[string]string{"error": "Content-Type must be application/json"})
			return
		}

		var req searchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var results []*vendors.VendorEntry
		for _, mac := range req.Macs {
			results = append(results, v.Search(mac)...)
		}

		json.NewEncoder(w).Encode(results)
	})
}
