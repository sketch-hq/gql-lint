package parser

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	gql "github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	gqlparser "github.com/vektah/gqlparser/v2/parser"
	gqlvalidator "github.com/vektah/gqlparser/v2/validator"
)

type SchemaField struct {
	Name string
}

type QueryField struct {
	Path              string
	SchemaPath        string
	IsDeprecated      bool
	DeprecationReason string
	File              string
	Line              int
}

type QueryFieldList []QueryField

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

func ParseSchemaFile(file string) (*ast.Schema, error) {
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

func ParseSchema(sources ...*ast.Source) (*ast.Schema, error) {
	schema, schemaerr := gqlvalidator.LoadSchema(sources...)
	if schemaerr != nil {
		return nil, schemaerr
	}

	return schema, nil
}

func ParseQuerySource(files []string, schema *ast.Schema) (QueryFieldList, error) {
	return queryTokensFromFiles(files, schema)
}

func ParseDeprecatedFields(schema *ast.Schema) []SchemaField {
	var fields []SchemaField

	for name, definition := range schema.Types {
		if ok, _ := isDeprecated(definition.Directives); ok {
			fields = append(fields, SchemaField{Name: name})
		}

		for _, field := range definition.Fields {
			if ok, _ := isDeprecated(field.Directives); ok {
				fields = append(fields, SchemaField{Name: fmt.Sprintf("%s.%s", name, field.Name)})
			}
		}
	}

	return fields
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

func queryTokensFromFiles(files []string, schema *ast.Schema) (QueryFieldList, error) {
	fields := QueryFieldList{}
	for _, file := range files {
		doc, err := parseQueryFile(file)
		if err != nil {
			return nil, err
		}
		// Ignoring errors here as there could be validation errors that we're not interested in.
		// We only call `.Validate` so the parser can populate the `Definition` fields.
		_ = gqlvalidator.Validate(schema, doc)
		fields = buildQueryTokens(doc, fields)
	}

	return fields, nil
}

func buildQueryTokens(query *ast.QueryDocument, fields QueryFieldList) QueryFieldList {
	for _, o := range query.Operations {
		fields = extractFields(o.SelectionSet, "", "", fields)
	}
	for _, f := range query.Fragments {
		fields = extractFields(f.SelectionSet, f.TypeCondition, "", fields)
	}

	return fields
}

func extractFields(set ast.SelectionSet, parentPath string, parentType string, fields QueryFieldList) QueryFieldList {
	for _, s := range set {
		switch f := s.(type) {
		case *ast.Field:
			path := ""
			if parentPath != "" {
				path = parentPath + "." + f.Name
			} else {
				path = f.Name
			}

			dep, depReason := isDeprecated(f.Definition.Directives)
			if !dep {
				fields = extractFields(f.SelectionSet, path, f.Definition.Type.Name(), fields)
				continue
			}

			field := QueryField{
				Path:              path,
				SchemaPath:        fmt.Sprintf("%s.%s", parentType, f.Name),
				File:              f.Position.Src.Name,
				Line:              f.Position.Line,
				IsDeprecated:      dep,
				DeprecationReason: depReason,
			}

			fields = append(fields, field)
			fields = extractFields(f.SelectionSet, path, f.Definition.Type.Name(), fields)
		}
	}
	return fields
}

func isDeprecated(dl ast.DirectiveList) (bool, string) {
	for _, d := range dl {
		if d.Name == "deprecated" {
			if reason, ok := d.ArgumentMap(make(map[string]interface{}))["reason"]; ok {
				return true, reason.(string)
			}
			return true, ""
		}
	}

	return false, ""
}
