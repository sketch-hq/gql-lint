package parser_test

import (
	"os"
	"testing"

	"github.com/matryer/is"
	"github.com/sketch-hq/gql-lint/parser"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/validator"
)

func TestParseSchema(t *testing.T) {
	is := is.New(t)

	schema, err := parser.ParseSchema(
		source(t, "testdata/schemas/with_deprecations.gql"),
		validator.Prelude,
	)

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

func TestParseQueries(t *testing.T) {
	t.Run("returns empty list if queries are not for the given schema", func(t *testing.T) {
		is := is.New(t)

		schema, err := parser.ParseSchema(
			source(t, "testdata/schemas/with_deprecations.gql"),
			validator.Prelude,
		)
		is.NoErr(err)

		fields, err := parser.ParseQueries(
			schema,
			source(t, "testdata/queries/for_other_schema.gql")
		)
		is.NoErr(err)

		is.Equal(len(fields), 0)
	})

	t.Run("successfully parses queries", func(t *testing.T) {
		is := is.New(t)

		schema, err := parser.ParseSchema(
			source(t, "testdata/schemas/with_deprecations.gql"),
			validator.Prelude,
		)
		is.NoErr(err)

		fields, err := parser.ParseQueries(
			schema,
			source(t, "testdata/queries/deprecation.gql"),
			source(t, "testdata/queries/one.gql"),
		)
		is.NoErr(err)

		is.Equal(len(fields), 1)

		field := fields[0]
		is.Equal(field.Path, "author.books.title")
		is.Equal(field.SchemaPath, "Book.title")
		is.True(field.IsDeprecated)
		is.Equal(field.File, "testdata/queries/deprecation.gql")
		is.Equal(field.Line, 7)
	})
}

func TestParseDeprecatedFields(t *testing.T) {
	is := is.New(t)

	schema, err := parser.ParseSchema(
		source(t, "testdata/schemas/with_deprecations.gql"),
		validator.Prelude,
	)
	is.NoErr(err)

	fields := parser.ParseDeprecatedFields(schema)

	is.Equal(len(fields), 1)

	field := fields[0]
	is.Equal(field.Name, "Book.title")
}

func source(t *testing.T, file string) *ast.Source {
	t.Helper()

	contents, err := os.ReadFile(file)
	if err != nil {
		t.Logf("Could not read test fixture file: %s", file)
		t.FailNow()
	}

	return &ast.Source{Name: file, Input: string(contents)}
}
