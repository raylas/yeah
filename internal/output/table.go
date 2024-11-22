package output

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type Table struct {
	w tabwriter.Writer
}

func NewTable(writer io.Writer) *Table {
	output := &Table{}
	output.w.Init(writer, 0, 0, 2, ' ', 0)
	return output
}

func (t *Table) WriteHeader(columns []string) error {
	painted := make([]string, len(columns))
	for c, column := range columns {
		painted[c] = fmt.Sprintf("%v%v%v", Default, column, Reset)
	}
	_, err := fmt.Fprintln(&t.w, strings.Join(painted, "\t"))
	return err
}

func (t *Table) WriteResource(fields []Field) error {
	painted := make([]string, len(fields))
	for f, field := range fields {
		color := Default
		if field.C != "" {
			color = field.C
		}
		painted[f] = fmt.Sprintf("%v%v%v", color, field.T, Reset)
	}
	_, err := fmt.Fprintln(&t.w, strings.Join(painted, "\t"))
	return err
}

func (t *Table) Flush() error {
	return t.w.Flush()
}
