package output

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type Field struct {
	Field             string `json:"field"`
	File              string `json:"file,omitempty"`
	Line              int    `json:"line"`
	DeprecationReason string `json:"reason,omitempty"`
}

func (f *Field) Equals(b *Field) bool {
	return f.Field == b.Field
}

type Data map[string][]Field

func (d Data) SortByField() {
	for schema := range d {
		sort.Slice(d[schema], func(i, j int) bool {
			return d[schema][i].Field < d[schema][j].Field
		})
	}
}

type DataWalkFunc func(schema string, field Field, fieldIdx int)

// Walk calls a walker function for each schema and field. Schemas are walked in alphabetical order.
func (d *Data) Walk(walker DataWalkFunc) {
	schemas := make([]string, 0, len(*d))
	for k := range *d {
		schemas = append(schemas, k)
	}
	sort.Strings(schemas)

	for _, schema := range schemas {
		fields := (*d)[schema]
		for fieldIdx, f := range fields {
			walker(schema, f, fieldIdx)
		}
	}
}

// AppendFields appends the given field to the list of fields for the given schema
func (d *Data) AppendField(schema string, field Field) {
	(*d)[schema] = append((*d)[schema], field)
}

// Diff compares this structure against `b` and returns a new Data instance with
// all the fields in `b` not found in `a`.
func (a Data) Diff(b Data) Data {
	result := Data{}
	// This is a dumb loop that could be optimized
	for schema, bFields := range b {
		for _, bField := range bFields {
			found := false
			if _, ok := a[schema]; ok {
				for _, aField := range a[schema] {
					if bField.Equals(&aField) {
						found = true
						break
					}
				}
			}

			if !(found) {
				result[schema] = append(result[schema], bField)
			}
		}
	}

	return result
}

type UnusedField struct {
	Field string `json:"field"`
}

func CompareFiles(fileA string, fileB string) (Data, error) {
	aContent, err := os.ReadFile(fileA)
	if err != nil {
		return nil, err
	}

	bContent, err := os.ReadFile(fileB)
	if err != nil {
		return nil, err
	}

	aData := Data{}
	err = json.Unmarshal(aContent, &aData)
	if err != nil {
		return nil, fmt.Errorf("could not json decode %s: %w", fileA, err)
	}

	bData := Data{}
	err = json.Unmarshal(bContent, &bData)
	if err != nil {
		return nil, fmt.Errorf("could not json decode %s: %w", fileB, err)
	}

	return aData.Diff(bData), nil
}
