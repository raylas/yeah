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
	C Color  // Color
}

type Color string

const (
	Reset   Color = "\x1b[0000m"
	Default Color = "\x1b[0039m"
	Green   Color = "\x1b[0032m"
	Red     Color = "\x1b[0031m"
)

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

func Colorize(color Color, format string, a ...interface{}) string {
	return fmt.Sprintf("%v%v%v", color, fmt.Sprintf(format, a...), Reset)
}

func Link(text, url string) string {
	return fmt.Sprintf("\x1b]8;;%v\x1b\\%v\x1b]8;;\x1b\\", url, text)
}
