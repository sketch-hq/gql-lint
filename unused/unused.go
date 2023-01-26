package unused

import (
	"github.com/sketch-hq/gql-lint/parser"
	"github.com/sketch-hq/gql-lint/sources"
	"github.com/vektah/gqlparser/v2/ast"
)

type UnusedField struct {
	Name string
}

type unusedRegistry map[string]UnusedField

func (r unusedRegistry) Record(field parser.SchemaField) {
	r[field.Name] = UnusedField{Name: field.Name}
}

func GetUnusedFields(schema *ast.Schema, queriesPaths []string) ([]UnusedField, error) {
	unusedFields := make(unusedRegistry)

	deprecatedFields := parser.ParseDeprecatedFields(schema)

	queries, err := sources.LoadQueries(schema, queriesPaths)
	if err != nil {
		return []UnusedField{}, err
	}

	for _, deprecatedField := range deprecatedFields {
		used := isUsed(deprecatedField, queries)
		recorded := isRecorded(deprecatedField, unusedFields)

		if used && recorded {
			delete(unusedFields, deprecatedField.Name)
		} else if !used && !recorded {
			unusedFields.Record(deprecatedField)
		}
	}

	result := make([]UnusedField, 0, len(unusedFields))
	for _, field := range unusedFields {
		result = append(result, field)
	}

	return result, nil
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
