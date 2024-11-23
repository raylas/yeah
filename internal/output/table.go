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
	_, err := fmt.Fprintln(&t.w, strings.Join(columns, "\t"))
	return err
}

func (t *Table) WriteResource(fields []Field) error {
	stringFields := make([]string, len(fields))
	for i, f := range fields {
		stringFields[i] = f.T
	}
	_, err := fmt.Fprintln(&t.w, strings.Join(stringFields, "\t"))
	return err
}

func (t *Table) Flush() error {
	return t.w.Flush()
}
