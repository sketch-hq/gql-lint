package parser

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	gql "github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	gqlparser "github.com/vektah/gqlparser/v2/parser"
)

type Field struct {
	Name              string
	Path              string
	Position          *ast.Position
	IsDeprecated      bool
	DeprecationReason string
}

type FieldKey string

type Fields map[FieldKey]Field

func ParseQueryDir(dir string) (Fields, error) {
	fields := make(Fields)
	files := findQueryFiles(dir)

	for _, file := range files {
		doc, err := parseQueryFile(file)
		if err != nil {
			return nil, err
		}
		buildQueryTokens(doc, fields)
	}

	return fields, nil
}

func ParseSchemaFile(file string) (Fields, error) {
	fields := make(Fields)
	schema, err := parseSchemaFile(file)
	if err != nil {
		return nil, err
	}
	buildSchemaTokens(schema, fields)
	return fields, nil
}

func parseSchemaFile(file string) (*ast.Schema, error) {
	contents, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	source := &ast.Source{Name: file, Input: string(contents)}
	schema, schemaerr := gql.LoadSchema(source)
	if schemaerr != nil {
		return nil, schemaerr
	}

	return schema, nil
}

func buildSchemaTokens(schema *ast.Schema, fields Fields) {
	for _, t := range schema.Types {
		for _, field := range t.Fields {
			dep, deprecationReason := isDeprecated(field.Directives)
			if !dep {
				continue
			}
			path := t.Name + "." + field.Name
			token := Field{
				Path:              path,
				IsDeprecated:      dep,
				DeprecationReason: deprecationReason,
			}
			fields[FieldKey(strings.ToLower(path))] = token
		}
	}
}

func findQueryFiles(startDir string) []string {
	files := []string{}
	filepath.WalkDir(startDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files
}

func parseQueryFile(file string) (*ast.QueryDocument, error) {
	queryContents, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	source := &ast.Source{Name: file, Input: string(queryContents)}
	doc, queryerr := gqlparser.ParseQuery(source)
	if queryerr != nil {
		return nil, queryerr
	}
	return doc, nil
}

func buildQueryTokens(query *ast.QueryDocument, fields Fields) {
	for _, o := range query.Operations {
		extractFields(o.SelectionSet, "", fields)
	}
	for _, f := range query.Fragments {
		extractFields(f.SelectionSet, f.TypeCondition, fields)
	}
}

func extractFields(set ast.SelectionSet, parentPath string, fields Fields) {
	for _, s := range set {
		switch f := s.(type) {
		case *ast.Field:
			path := ""
			if parentPath != "" {
				path = parentPath + "." + f.Name
			} else {
				path = f.Name
			}
			field := Field{
				Name:     f.Name,
				Path:     path,
				Position: f.Position,
			}

			fields[FieldKey(strings.ToLower(path))] = field
			extractFields(f.SelectionSet, path, fields)
		}
	}
}

func isDeprecated(dl ast.DirectiveList) (bool, string) {
	for _, d := range dl {
		if d.Name == "deprecated" {
			return true, ""
		}
	}

	return false, ""
}
