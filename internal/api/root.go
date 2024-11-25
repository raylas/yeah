package api

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"rymnd.net/yeah/internal/vendors"
)

//go:embed templates/root.html
var rootHTML string

func handleHtmlRoot(v *vendors.Vendors) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("root").Parse(rootHTML)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, v)
	})
}

func handleRoot(v *vendors.Vendors) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")
		if isCurl(userAgent) {
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "Available routes:\n/<mac>[,<mac>...]   Search vendors\n/sources            List sources\n")
			return
		}
		handleHtmlRoot(v).ServeHTTP(w, r)
	})
}

func isCurl(userAgent string) bool {
	return strings.Contains(strings.ToLower(userAgent), "curl")
}
