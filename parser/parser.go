package parser

import (
	"io/fs"
	"os"
	"path/filepath"

	gql "github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	gqlparser "github.com/vektah/gqlparser/v2/parser"
	gqlvalidator "github.com/vektah/gqlparser/v2/validator"
)

type QueryField struct {
	Path              string
	IsDeprecated      bool
	DeprecationReason string
	File              string
	Line              int
}

type QueryFieldList []QueryField

func ParseQueryDir(dir string, schema *ast.Schema) (QueryFieldList, error) {
	fields := QueryFieldList{}
	files := findQueryFiles(dir)
	// @todo: error out if no files are found

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

func ParseSchemaFile(file string) (*ast.Schema, error) {
	return parseSchemaFile(file)
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

func buildQueryTokens(query *ast.QueryDocument, fields QueryFieldList) QueryFieldList {
	for _, o := range query.Operations {
		fields = extractFields(o.SelectionSet, "", fields)
	}
	for _, f := range query.Fragments {
		fields = extractFields(f.SelectionSet, f.TypeCondition, fields)
	}

	return fields
}

func extractFields(set ast.SelectionSet, parentPath string, fields QueryFieldList) QueryFieldList {
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
				fields = extractFields(f.SelectionSet, path, fields)
				continue
			}

			field := QueryField{
				Path:              path,
				File:              f.Position.Src.Name,
				Line:              f.Position.Line,
				IsDeprecated:      dep,
				DeprecationReason: depReason,
			}

			fields = append(fields, field)
			fields = extractFields(f.SelectionSet, path, fields)
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
