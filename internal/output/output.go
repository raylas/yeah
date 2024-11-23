package output

import (
	"fmt"
	"os"
)

type Writer interface {
	WriteHeader(columns []string) error
	WriteResource(fields []Field) error
	Flush() error
}

type Field struct {
	T string // Text
}

func NewWriter(format string) (Writer, error) {
	writer := os.Stdout
	switch {
	case format == "table":
		return NewTable(writer), nil
	case format == "json":
		return NewJSON(writer), nil
	default:
		return nil, fmt.Errorf("output format %q is not supported", format)
	}
}
