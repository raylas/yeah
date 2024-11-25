package output

import (
	"fmt"
	"io"
	"strings"
)

type Writer interface {
	WriteHeader(columns []string) error
	WriteResource(fields []Field) error
	Flush() error
}

type Field struct {
	T string // Text
}

func NewWriter(writer io.Writer, format string) (Writer, error) {
	switch {
	case format == "table":
		return NewTable(writer), nil
	case format == "json":
		return NewJSON(writer), nil
	case format == "html":
		return NewHTML(writer), nil
	default:
		return nil, fmt.Errorf("output format %q is not supported", format)
	}
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(strings.Replace(fmt.Sprint(v), "\n", "", -1))
}
