package otterbeat

import (
	"reflect"
	"strings"
)

// StructCreator is a struct for the golang struct creator.
type StructCreator struct {
	SlicePtr       reflect.Value
	SliceValue     reflect.Value
	ItemStructType reflect.Type
	Columns        []itemStructField
	ColumnsMap     map[int]itemStructField
}

// PrepareColumns prepares the struct columns relation by the column titles.
func (c *StructCreator) PrepareColumns(columnTitles []string) {
	m := make(map[int]itemStructField)

	for _, structCol := range c.Columns {
		for j, title := range columnTitles {
			if strings.Contains(structCol.ColumnName, title) {
				m[j] = structCol
				break
			}
		}
	}

	c.ColumnsMap = m
}

// CreateSliceItem creates a slice item from the columns values.
func (c *StructCreator) CreateSliceItem(columns []string) {
	v := reflect.New(c.ItemStructType).Elem()

	for j, col := range columns {
		if structCol, ok := c.ColumnsMap[j]; ok {
			v.Field(structCol.FiledIndex).Set(reflect.ValueOf(col))
		}
	}

	c.SliceValue = reflect.Append(c.SliceValue, v)
	c.SlicePtr.Elem().Set(c.SliceValue)
}

type itemStructField struct {
	FiledIndex int
	ColumnName string
}

// NewStructCreator creates a new StructCreator.
func NewStructCreator(slicePtr interface{}) *StructCreator {
	SlicePtr := reflect.ValueOf(slicePtr)
	sliceValue := SlicePtr.Elem()

	s := &StructCreator{
		SlicePtr:       SlicePtr,
		SliceValue:     sliceValue,
		ItemStructType: sliceValue.Type().Elem(),
	}

	columns := make([]itemStructField, 0)

	for i := 0; i < s.ItemStructType.NumField(); i++ {
		fi := s.ItemStructType.Field(i)
		if fi.PkgPath != "" {
			continue
		}

		colName := fi.Tag.Get("col")
		if colName == "" {
			colName = fi.Name
		}

		columns = append(columns, itemStructField{FiledIndex: i, ColumnName: colName})
	}

	s.Columns = columns

	return s
}
