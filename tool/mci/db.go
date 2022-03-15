package mci

import (
	"fmt"
	"io"

	"github.com/bingoohuang/sqlx"
	"github.com/jedib0t/go-pretty/v6/table"
)

// PrintSQLResult prints the result r of sqlStr execution.
func PrintSQLResult(stdout, stderr io.Writer, sqlStr string, r sqlx.ExecResult) error {
	if r.Error != nil {
		fmt.Fprintf(stderr, "error %v\n", r.Error)
		return r.Error
	}

	fmt.Fprintf(stdout, "SQL: %s\n", sqlStr)
	fmt.Fprintf(stdout, "Cost: %s\n", r.CostTime.String())

	if !r.IsQuerySQL {
		return nil
	}

	cols := len(r.Headers) + 1
	header := make(table.Row, cols)
	header[0] = "#"

	for i, h := range r.Headers {
		header[i+1] = h
	}

	t := table.NewWriter()
	t.SetOutputMirror(stdout)
	t.AppendHeader(header)

	for i, r := range r.Rows {
		row := make(table.Row, cols)
		row[0] = i + 1

		for j, c := range r {
			row[j+1] = c
		}

		t.AppendRow(row)
	}

	t.Render()

	return nil
}
