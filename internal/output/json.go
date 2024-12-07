package output

import (
	"encoding/json"
	"io"
	"strings"
)

type JSON struct {
	w io.Writer
	h []string
}

func NewJSON(writer io.Writer) *JSON {
	return &JSON{
		w: writer,
	}
}

func (j *JSON) WriteHeader(columns []string) error {
	headers := make([]string, len(columns))
	for c, column := range columns {
		column = strings.Replace(column, " ", "_", -1)
		headers[c] = strings.ToLower(column)
	}
	j.h = headers
	return nil
}

func (j *JSON) WriteResource(fields []Field) error {
	data := make(map[string]interface{})
	for f, field := range fields {
		data[j.h[f]] = field.T
	}
	return writeAsIndentedJSON(j.w, data)
}

func writeAsIndentedJSON(wr io.Writer, data interface{}) error {
	enc := json.NewEncoder(wr)
	enc.SetIndent("", "    ")
	return enc.Encode(data)
}

func (j *JSON) Flush() error {
	return nil
}
