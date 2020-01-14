package mci

import (
	"fmt"
	"os"
	"reflect"
	"regexp"

	"github.com/bingoohuang/gou/reflec"
	"github.com/jedib0t/go-pretty/table"
)

// TablePrinter print table
type TablePrinter struct {
	dittoMark string
}

// Print prints the table
func (p TablePrinter) Print(value interface{}) {
	header := make(table.Row, 0)
	rows := make([]table.Row, 0)

	header = append(header, "#")

	v := reflect.ValueOf(value)
	if v.IsNil() {
		return
	}

	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
		fallthrough
	case reflect.Struct:
		fields := reflec.CachedStructFields(v.Type(), "header")
		createHeader(fields, &header)
		createRow(fields, 0, v, &rows)
	case reflect.Slice:
		if v.Len() == 0 {
			return
		}

		fields := reflec.CachedStructFields(v.Type().Elem(), "header")
		createHeader(fields, &header)

		for i := 0; i < v.Len(); i++ {
			createRow(fields, i, v.Index(i), &rows)
		}
	default:
		return
	}

	p.tableRender(header, rows...)
	fmt.Println()
}

func createRow(fields []reflec.StructField, rowIndex int, v reflect.Value, rows *[]table.Row) {
	row := make(table.Row, 0)
	row = append(row, rowIndex+1) // nolint gomnd

	for _, f := range fields {
		row = append(row, v.Field(f.Index).Interface())
	}

	*rows = append(*rows, row)
}

func createHeader(fields []reflec.StructField, header *table.Row) {
	for _, f := range fields {
		*header = append(*header, BlankCamel(f.Name))
	}
}

func (p TablePrinter) tableRender(header table.Row, rows ...table.Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(header)

	if p.dittoMark != "" {
		t.AppendRows(p.dittoMarkRows(rows))
	} else {
		t.AppendRows(rows)
	}

	t.Render()
}

func (p TablePrinter) dittoMarkRows(rows []table.Row) []table.Row {
	mark := make(map[int]interface{})

	for i, row := range rows {
		for j, cell := range row {
			v, ok := mark[j]
			if ok && v != "" && v == cell {
				rows[i][j] = p.dittoMark
			} else {
				mark[j] = cell
			}
		}
	}

	return rows
}

// BlankCamel make a camel string to blanks
func BlankCamel(str string) string {
	blank := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(str, "${1} ${2}")
	return regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(blank, "${1} ${2}")
}
