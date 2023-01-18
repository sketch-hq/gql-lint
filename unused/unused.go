package unused

import "github.com/sketch-hq/gql-lint/parser"

func GetUnusedFields(schemaPath string, queriesPath string) ([]parser.SchemaField, error) {
	var unusedFields []parser.SchemaField

	schema, err := parser.ParseSchemaFile(schemaPath)
	if err != nil {
		return unusedFields, err
	}

	deprecatedFields := parser.ParseDeprecatedFields(schema)

	queries, err := parser.ParseQuerySource(queriesPath, schema)
	if err != nil {
		return unusedFields, err
	}

	for _, deprecatedField := range deprecatedFields {
		if !IsFieldUsed(deprecatedField, queries) {
			unusedFields = append(unusedFields, deprecatedField)
		}
	}

	return unusedFields, nil
}

func IsFieldUsed(field parser.SchemaField, queryFields parser.QueryFieldList) bool {
	for _, queryField := range queryFields {
		if queryField.SchemaPath == field.Name {
			return true
		}
	}
	return false
}
