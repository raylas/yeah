package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"rymnd.net/yeah/internal/output"
	"rymnd.net/yeah/internal/vendors"
)

func handleRawSearch(v *vendors.Vendors) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer("").Start(r.Context(), "handleRawSearch")
		defer span.End()
		setCommonAttributes(ctx, r)

		query := r.PathValue("query")
		macs := strings.Split(query, ",")

		results := make([]*vendors.VendorEntry, 0)
		for _, mac := range macs {
			results = append(results, v.Search(ctx, mac)...)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func handleHtmlSearch(v *vendors.Vendors) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer("").Start(r.Context(), "handleHtmlSearch")
		defer span.End()
		setCommonAttributes(ctx, r)

		query := r.PathValue("query")
		macs := strings.Split(query, ",")

		results := make([]*vendors.VendorEntry, 0)
		for _, mac := range macs {
			results = append(results, v.Search(ctx, mac)...)
		}

		b := &bytes.Buffer{}
		o, err := output.NewWriter(b, "html")
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		o.WriteHeader([]string{"OUI", "Name"})
		for _, result := range results {
			o.WriteResource([]output.Field{{T: result.Oui}, {T: result.Name}})
		}
		o.Flush()

		w.Header().Set("Content-Type", "text/html")
		w.Write(b.Bytes())
	})
}

func handleSources(v *vendors.Vendors) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer("").Start(r.Context(), "handleSources")
		defer span.End()
		setCommonAttributes(ctx, r)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(v.Sources); err != nil {
			span.SetStatus(codes.Error, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
