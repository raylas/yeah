package output

import (
	"io"
	"strings"
)

type HTML struct {
	w             io.Writer
	h             []string
	headerWritten bool
}

func NewHTML(writer io.Writer) *HTML {
	html := &HTML{w: writer}
	io.WriteString(writer, "<table>\n")
	return html
}

func (h *HTML) WriteHeader(columns []string) error {
	io.WriteString(h.w, "  <thead>\n    <tr>\n")
	for _, column := range columns {
		io.WriteString(h.w, `      <th scope="col">`)
		io.WriteString(h.w, column)
		io.WriteString(h.w, "</th>\n")
	}
	io.WriteString(h.w, "    </tr>\n  </thead>\n  <tbody>\n")

	headers := make([]string, len(columns))
	for i, column := range columns {
		column = strings.Replace(column, " ", "_", -1)
		headers[i] = strings.ToLower(column)
	}
	h.h = headers
	h.headerWritten = true
	return nil
}

func (h *HTML) WriteResource(fields []Field) error {
	io.WriteString(h.w, "    <tr>\n")
	for _, field := range fields {
		io.WriteString(h.w, "      <td>")
		io.WriteString(h.w, toString(field.T))
		io.WriteString(h.w, "</td>\n")
	}
	io.WriteString(h.w, "    </tr>\n")
	return nil
}

func (h *HTML) Flush() error {
	io.WriteString(h.w, "  </tbody>\n</table>\n")
	return nil
}
