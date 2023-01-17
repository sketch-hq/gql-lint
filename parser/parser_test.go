package parser

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestParseSchemaFile(t *testing.T) {
	is := is.New(t)

	schema, err := ParseSchemaFile("testdata/schemas/with_deprecations.gql")

	is.NoErr(err)

	deprecatedField := schema.Types["Book"].Fields[1]

	is.Equal(deprecatedField.Name, "title")

	deprecated := false
	for _, directive := range deprecatedField.Directives {
		if directive.Name == "deprecated" {
			is.Equal(directive.Arguments[0].Name, "reason")
			is.Equal(directive.Arguments[0].Value.String(), `"untitled books are better"`)

			deprecated = true
		}
	}

	is.True(deprecated)

	field := schema.Types["Book"].Fields[2]
	is.Equal(field.Name, "author")

	deprecated = false
	for _, directive := range field.Directives {
		if directive.Name == "deprecated" {
			deprecated = true
		}
	}
	is.True(!deprecated)
}

func TestParseSchemaFile_NotFound(t *testing.T) {
	is := is.New(t)

	_, err := ParseSchemaFile("testdata/schemas/not_found.gql")

	is.True(strings.Contains(err.Error(), "open testdata/schemas/not_found.gql: no such file or directory"))
}

func TestParseQueryDir(t *testing.T) {
	is := is.New(t)

	schema, err := ParseSchemaFile("testdata/schemas/with_deprecations.gql")
	is.NoErr(err)

	fields, err := ParseQueryDir("testdata/queries", schema)
	is.NoErr(err)

	is.Equal(len(fields), 1)

	field := fields[0]
	is.Equal(field.Path, "author.books.title")
	is.Equal(field.SchemaPath, "Book.title")
	is.True(field.IsDeprecated)
	is.Equal(field.File, "testdata/queries/deprecation.gql")
	is.Equal(field.Line, 7)
}

func TestParseDeprecatedFields(t *testing.T) {
	is := is.New(t)

	schema, err := ParseSchemaFile("testdata/schemas/with_deprecations.gql")
	is.NoErr(err)

	fields := ParseDeprecatedFields(schema)

	is.Equal(len(fields), 1)

	field := fields[0]
	is.Equal(field.Name, "Book.title")
}

func TestIsFieldUsed(t *testing.T) {
	is := is.New(t)

	schema, err := ParseSchemaFile("testdata/schemas/with_deprecations.gql")
	is.NoErr(err)

	query, err := ParseQueryDir("testdata/queries/deprecation.gql", schema)
	is.NoErr(err)

	used := IsFieldUsed(&SchemaField{Name: "Book.title"}, query)
	is.True(used)

	query, err = ParseQueryDir("testdata/queries/one.gql", schema)
	is.NoErr(err)

	used = IsFieldUsed(&SchemaField{Name: "Book.title"}, query)
	is.True(!used)
}
