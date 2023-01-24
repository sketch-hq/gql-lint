package schema

import (
	"os"
	"strings"

	"github.com/sketch-hq/gql-lint/introspection"
	"github.com/sketch-hq/gql-lint/parser"
	"github.com/vektah/gqlparser/v2/ast"
	gqlvalidator "github.com/vektah/gqlparser/v2/validator"
)

const deprecated = `directive @deprecated(
	reason: String = "No longer supported"
) on FIELD_DEFINITION | ARGUMENT_DEFINITION | ENUM_VALUE | INPUT_FIELD_DEFINITION`

func Load(source string) (*ast.Schema, error) {
	loader := loaderForSource(source)
	sources, err := loader.Load(source)
	if err != nil {
		return nil, err
	}

	return parser.ParseSchema(sources...)
}

type Loader interface {
	Load(source string) ([]*ast.Source, error)
}

type FileLoader struct{}

func (s FileLoader) Load(source string) ([]*ast.Source, error) {
	contents, err := os.ReadFile(source)
	if err != nil {
		return nil, err
	}

	return []*ast.Source{
		gqlvalidator.Prelude,
		{Name: source, Input: string(contents), BuiltIn: false},
	}, nil
}

type HttpLoader struct{}

func (s HttpLoader) Load(source string) ([]*ast.Source, error) {
	contents, err := introspection.Load(source)
	if err != nil {
		return nil, err
	}

	schemaString := string(contents)
	sources := []*ast.Source{{Name: source, Input: schemaString, BuiltIn: true}}

	// workaround for absinthe based graphql servers that does NOT include the deprecated
	// directive in the schema.
	if !strings.Contains(schemaString, "directive @deprecated") {
		sources = append(sources, &ast.Source{Input: deprecated, BuiltIn: false})
	}

	return sources, nil
}

func loaderForSource(source string) Loader {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return HttpLoader{}
	}
	return FileLoader{}
}
