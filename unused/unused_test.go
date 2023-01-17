package unused

import (
	"testing"

	"github.com/matryer/is"
	"github.com/sketch-hq/gql-lint/parser"
)

func TestIsFieldUsed(t *testing.T) {
	is := is.New(t)

	schema, err := parser.ParseSchemaFile("testdata/schemas/with_deprecations.gql")
	is.NoErr(err)

	query, err := parser.ParseQueryDir("testdata/queries/deprecation.gql", schema)
	is.NoErr(err)

	used := IsFieldUsed(&parser.SchemaField{Name: "Book.title"}, query)
	is.True(used)

	query, err = parser.ParseQueryDir("testdata/queries/one.gql", schema)
	is.NoErr(err)

	used = IsFieldUsed(&parser.SchemaField{Name: "Book.title"}, query)
	is.True(!used)
}
