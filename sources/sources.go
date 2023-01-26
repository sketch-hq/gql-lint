package sources

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/sketch-hq/gql-lint/introspection"
	"github.com/sketch-hq/gql-lint/parser"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/validator"
)

const deprecated = `directive @deprecated(
	reason: String = "No longer supported"
) on FIELD_DEFINITION | ARGUMENT_DEFINITION | ENUM_VALUE | INPUT_FIELD_DEFINITION`

// LoadSchema will load a schema from either a given file or url
func LoadSchema(source string) (*ast.Schema, error) {
	loader, err := loaderForSource(source, true)
	if err != nil {
		return nil, err
	}
	sources, err := loader.Load(source)
	if err != nil {
		return nil, err
	}

	return parser.ParseSchema(sources...)
}

// LoadQueries will load one or more queries from either a file or a directory
func LoadQueries(schema *ast.Schema, sources []string) (parser.QueryFieldList, error) {
	allSources := []*ast.Source{}
	for _, source := range sources {
		loader, err := loaderForSource(source, false)
		if err != nil {
			return nil, err
		}

		sources, err := loader.Load(source)
		if err != nil {
			return nil, err
		}

		allSources = append(allSources, sources...)
	}

	return parser.ParseQueries(schema, allSources...)
}

type Loader interface {
	Load(source string) ([]*ast.Source, error)
}

type fileLoader struct {
	IncludePrelude bool
}

func (s fileLoader) Load(source string) ([]*ast.Source, error) {
	contents, err := os.ReadFile(source)
	if err != nil {
		return nil, err
	}

	sources := []*ast.Source{
		{Name: source, Input: string(contents), BuiltIn: false},
	}

	if s.IncludePrelude {
		sources = append(sources, validator.Prelude)
	}
	return sources, nil

}

type directoryLoader struct{}

func (l directoryLoader) Load(source string) ([]*ast.Source, error) {
	var sources []*ast.Source

	filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		sources = append(sources, &ast.Source{Name: path, Input: string(contents)})
		return nil
	})

	return sources, nil
}

type httpLoader struct{}

func (s httpLoader) Load(source string) ([]*ast.Source, error) {
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

func loaderForSource(source string, prelude bool) (Loader, error) {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return httpLoader{}, nil
	}

	isDir, err := isDirectory(source)
	if err != nil {
		return nil, err
	}
	if isDir {
		return directoryLoader{}, nil
	}

	return fileLoader{IncludePrelude: prelude}, nil
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}
