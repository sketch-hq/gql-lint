package unused

import "github.com/sketch-hq/gql-lint/parser"

func IsFieldUsed(field *parser.SchemaField, queryFields parser.QueryFieldList) bool {
	for _, queryField := range queryFields {
		if queryField.SchemaPath == field.Name {
			return true
		}
	}
	return false
}
