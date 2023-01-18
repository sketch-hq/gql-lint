package unused

import (
	"testing"

	"github.com/matryer/is"
	"github.com/sketch-hq/gql-lint/parser"
)

func TestGetUnusedFields(t *testing.T) {
	is := is.New(t)

	unusedFields, err := GetUnusedFields("testdata/schemas/with_deprecations.gql", "testdata/queries/deprecation.gql")
	is.NoErr(err)

	is.Equal(len(unusedFields), 1)

	is.Equal(unusedFields[0].Name, "Book.oldTitle")
}

func TestIsFieldUsed(t *testing.T) {
	is := is.New(t)

	schema, err := parser.ParseSchemaFile("testdata/schemas/with_deprecations.gql")
	is.NoErr(err)

	query, err := parser.ParseQueryDir("testdata/queries/deprecation.gql", schema)
	is.NoErr(err)

	used := IsFieldUsed(parser.SchemaField{Name: "Book.title"}, query)
	is.True(used)

	query, err = parser.ParseQueryDir("testdata/queries/one.gql", schema)
	is.NoErr(err)

	used = IsFieldUsed(parser.SchemaField{Name: "Book.title"}, query)
	is.True(!used)
}
