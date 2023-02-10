package unused

import (
	"fmt"

	"github.com/sketch-hq/gql-lint/output"
	"github.com/sketch-hq/gql-lint/parser"
	"github.com/sketch-hq/gql-lint/sources"
	"github.com/vektah/gqlparser/v2/ast"
)

type UnusedField struct {
	Name string
}

type unusedRegistry map[string]output.Field

func (r unusedRegistry) Record(field parser.SchemaField) {
	r[field.Name] = output.Field{
		Field: field.Name,
		Line:  field.Line,
		// Not recording the file as that would be the same for all (the schema)
	}
}

func GetUnusedFields(schemas []string, queryFiles []string, verbose bool) (output.Data, error) {
	out := output.Data{}

	for _, schemaSource := range schemas {
		schema, err := sources.LoadSchema(schemaSource)
		if err != nil {
			return nil, err
		}

		if verbose {
			fmt.Println("debug: Succesfully loaded schema from", schemaSource)
		}

		queries, err := sources.LoadQueries(schema, queryFiles)
		if err != nil {
			return nil, err
		}

		out[schemaSource] = []output.Field{}
		unusedFields, err := getUnusedFields(schema, queries)
		if err != nil {
			return nil, err
		}
		for _, field := range unusedFields {
			out.AppendField(schemaSource, field)
		}
	}

	out.SortByField()

	return out, nil
}

func getUnusedFields(schema *ast.Schema, queries parser.QueryFieldList) (unusedRegistry, error) {
	unusedFields := make(unusedRegistry)
	deprecatedFields := parser.ParseDeprecatedFields(schema)

	for _, deprecatedField := range deprecatedFields {
		used := isUsed(deprecatedField, queries)
		recorded := isRecorded(deprecatedField, unusedFields)

		if used && recorded {
			delete(unusedFields, deprecatedField.Name)
		} else if !used && !recorded {
			unusedFields.Record(deprecatedField)
		}
	}

	return unusedFields, nil
}

func isUsed(field parser.SchemaField, queryFields parser.QueryFieldList) bool {
	for _, queryField := range queryFields {
		if queryField.SchemaPath == field.Name {
			return true
		}
	}
	return false
}

func isRecorded(field parser.SchemaField, registry unusedRegistry) bool {
	_, found := registry[field.Name]

	return found
}
