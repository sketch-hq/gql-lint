package schema

import (
	"os"
	"strings"

	"github.com/sketch-hq/gql-lint/introspection"
	"github.com/sketch-hq/gql-lint/parser"
	"github.com/vektah/gqlparser/v2/ast"
)

func Load(source string) (*ast.Schema, error) {
	var builtin bool
	var contents []byte
	var err error

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		builtin = true
		contents, err = introspection.Load(source)
	} else {
		contents, err = os.ReadFile(source)
	}

	if err != nil {
		return nil, err
	}

	return parser.ParseSchema(source, string(contents), builtin)
}
