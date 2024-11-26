package api

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"rymnd.net/yeah/internal/vendors"
)

//go:embed templates/root.html
var rootHTML string

func handleHtmlRoot(v *vendors.Vendors) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, span := otel.Tracer("").Start(r.Context(), "handleHtmlRoot")
		defer span.End()

		tmpl, err := template.New("root").Parse(rootHTML)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, v)
	})
}

func handleRoot(v *vendors.Vendors) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer("").Start(r.Context(), "handleRoot")
		defer span.End()
		setCommonAttributes(ctx, r)

		userAgent := r.Header.Get("User-Agent")
		if isCurl(userAgent) {
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "Available routes:\n/<mac>[,<mac>...]   Search vendors\n/sources            List sources\n")
			return
		}
		handleHtmlRoot(v).ServeHTTP(w, r.WithContext(ctx))
	})
}

func isCurl(userAgent string) bool {
	return strings.Contains(strings.ToLower(userAgent), "curl")
}
